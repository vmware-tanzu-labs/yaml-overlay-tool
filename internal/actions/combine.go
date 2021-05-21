// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

var ErrCombineOnlyForScalars = errors.New("combine action can only be used on scalar values")

func CombineNode(originalValue, newValue *yaml.Node) error {
	if originalValue.Kind != yaml.ScalarNode || newValue.Kind != yaml.ScalarNode {
		return ErrCombineOnlyForScalars
	}

	mergeComments(originalValue, newValue)

	switch originalValue.Tag {
	case "!!int":
		o, err := strconv.Atoi(originalValue.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to int for merging: %w", originalValue.Value, err)
		}

		n, err := strconv.Atoi(newValue.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to int for merging: %w", newValue.Value, err)
		}

		originalValue.Value = strconv.Itoa(o + n)

		return nil

	case "!!bool":
		o, err := strconv.ParseBool(originalValue.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to bool for merging: %w", originalValue.Value, err)
		}

		n, err := strconv.ParseBool(newValue.Value)
		if err != nil {
			return fmt.Errorf("failed to %s to bool for merging: %w", newValue.Value, err)
		}

		originalValue.Value = strconv.FormatBool(o && n)

		return nil
	}

	originalValue.Value += newValue.Value

	return nil
}
