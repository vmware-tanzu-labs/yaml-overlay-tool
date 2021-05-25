// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"

	"github.com/qdm12/reprint"
	"gopkg.in/yaml.v3"
)

func ReplaceNode(original, replace *yaml.Node) error {
	ov := *original

	if err := reprint.FromTo(replace, original); err != nil {
		return fmt.Errorf("failed to replace value: %w", err)
	}

	mergeComments(&ov, original)
	original.HeadComment = ov.HeadComment
	original.LineComment = ov.LineComment
	original.FootComment = ov.FootComment

	return nil
}
