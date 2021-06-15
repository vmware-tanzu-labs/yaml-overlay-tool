// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"strings"
)

const (
	// ValueMarker is a format marker for the previous value.
	ValueMarker = `%v`
	// LineCommentMarker is a format marker for the previous line comment.
	LineCommentMarker = `%l`
	// HeadCommentMarker is a format marker for the previous head comment.
	HeadCommentMarker = `%h`
	// FootCommentMarker is a format marker for the previous foot comment.
	FootCommentMarker = `%f`
	// KeyMarker is a format marker for the previous key name.
	KeyMarker = `%k`
)

func checkForMarkers(s string) bool {
	for _, m := range []string{
		ValueMarker,
		LineCommentMarker,
		HeadCommentMarker,
		FootCommentMarker,
		KeyMarker,
	} {
		if strings.Contains(s, m) {
			return true
		}
	}

	return false
}
