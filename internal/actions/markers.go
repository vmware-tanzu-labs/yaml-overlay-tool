// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"
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

func sanitizeMarkers(formatStr string, v ...interface{}) (sanitizedfmt string, sanitizedV []interface{}) {
	for i, t := range []string{
		ValueMarker,
		LineCommentMarker,
		HeadCommentMarker,
		FootCommentMarker,
		KeyMarker,
	} {
		for j := range v {
			v[j] = strings.ReplaceAll(v[j].(string), t, "")
		}

		formatStr = strings.ReplaceAll(formatStr, t, fmt.Sprintf("%%[%d]v", i+1))
	}

	return formatStr, v
}

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
