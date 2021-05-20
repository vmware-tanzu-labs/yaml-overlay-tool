// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import "strings"

type multiString []string

func (ms *multiString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err == nil {
		*ms = []string{s}

		return nil
	}

	type ss []string

	return unmarshal((*ss)(ms))
}

func (ms multiString) String() string {
	return strings.Join(ms, ",")
}
