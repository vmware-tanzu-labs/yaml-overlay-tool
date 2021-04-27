package actions

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

var ErrMergeMustBeOfSameKind = errors.New("both values must be of same type to use 'merge' action")

func Merge(originalValue, newValue *yaml.Node) error {
	if originalValue.Kind != newValue.Kind {
		// are both originalValue and newValue the same 'Kind'?
		return ErrMergeMustBeOfSameKind
	}

	switch originalValue.Kind {
	case yaml.ScalarNode:
		// scalar:
		//   orig + new
		if err := appendValues(originalValue, newValue); err != nil {
			return err
		}
	case yaml.SequenceNode:
		// sequence:
		//   originalValue extended with newValue
		originalValue.Content = append(originalValue.Content, newValue.Content...)

		return nil
	case yaml.MappingNode, yaml.DocumentNode, yaml.AliasNode:
		// maps:
		//	 recursive merge of data
		if err := mergo.Merge(originalValue, *newValue, mergo.WithOverride); err != nil {
			return fmt.Errorf("an error occurred during merge: %w", err)
		}
	}

	return nil
}

func appendValues(ov, nv *yaml.Node) error {
	switch ov.Tag {
	case "!!int":
		o, err := strconv.Atoi(ov.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to int for merging: %w", ov.Value, err)
		}

		n, err := strconv.Atoi(nv.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to int for merging: %w", nv.Value, err)
		}

		ov.Value = strconv.Itoa(o + n)

	case "!!bool":
		o, err := strconv.ParseBool(ov.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to bool for merging: %w", ov.Value, err)
		}

		n, err := strconv.ParseBool(nv.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to bool for merging: %w", nv.Value, err)
		}

		ov.Value = strconv.FormatBool(o && n)

	default:
		ov.Value += ov.Value
	}

	return nil
}
