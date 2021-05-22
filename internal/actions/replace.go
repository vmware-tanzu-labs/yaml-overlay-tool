// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"

	"github.com/ulule/deepcopier"
	"gopkg.in/yaml.v3"
)

func ReplaceNode(original, replaceValue *yaml.Node) error {
	ov := *original

	mergeComments(&ov, replaceValue)

	if err := deepcopier.Copy(replaceValue).To(original); err != nil {
		return fmt.Errorf("failed to replace value: %w", err)
	}

	original.HeadComment = ov.HeadComment
	original.LineComment = ov.LineComment
	original.FootComment = ov.FootComment

	return nil
}
