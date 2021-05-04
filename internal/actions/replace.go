// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"

	"github.com/jinzhu/copier"
	"gopkg.in/yaml.v3"
)

func Replace(original, replaceValue *yaml.Node) error {
	options := copier.Option{
		IgnoreEmpty: false,
		DeepCopy:    true,
	}

	if err := copier.CopyWithOption(original, replaceValue, options); err != nil {
		return fmt.Errorf("failed to replace value: %w", err)
	}

	return nil
}
