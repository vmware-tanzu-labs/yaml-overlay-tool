// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

var ErrFormatOnlyForScalars = errors.New("format action can only be used on scalar values")

func FormatNode(originalValue, newValue *yaml.Node) error {
	if originalValue.Kind == yaml.ScalarNode && newValue.Kind == yaml.ScalarNode {
		ov := originalValue.Value

		if err := mergo.Merge(originalValue, *newValue, mergo.WithOverride); err != nil {
			return fmt.Errorf("failed to merge prior to formatting: %w", err)
		}

		originalValue.Value = CondSprintf(originalValue.Value, ov)

		return nil
	}

	return ErrFormatOnlyForScalars
}

func CondSprintf(format string, v ...interface{}) string {
	v = append(v, "")
	format += fmt.Sprint("%[", len(v), "]s")

	return fmt.Sprintf(format, v...)
}
