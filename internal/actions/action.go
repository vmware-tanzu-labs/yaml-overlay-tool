// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidAction = errors.New("invalid overlay action")

type Action int

const (
	Invalid = iota
	Merge
	Replace
	Delete
	Format
	Math
)

func (a Action) String() string {
	toString := map[Action]string{
		Invalid: "",
		Merge:   "merge",
		Replace: "replace",
		Delete:  "delete",
		Math:    "math",
	}

	return toString[a]
}

func (a *Action) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var y string

	if err := unmarshal(&y); err != nil {
		return err
	}

	y = strings.ToLower(y)

	toID := map[string]Action{
		"":        Invalid,
		"merge":   Merge,
		"replace": Replace,
		"delete":  Delete,
		"math":    Math,
	}

	if toID[y] == Invalid {
		return fmt.Errorf("%w, %s", ErrInvalidAction, y)
	}

	*a = toID[y]

	return nil
}

func (a Action) MarshalYAML() (interface{}, error) {
	return a.String(), nil
}

type OnMissingAction int

const (
	Ignore = iota
	Inject
)

func (a OnMissingAction) String() string {
	toString := map[OnMissingAction]string{
		Ignore: "ignore",
		Inject: "inject",
	}

	return toString[a]
}

func (a *OnMissingAction) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var y string

	if err := unmarshal(&y); err != nil {
		return err
	}

	y = strings.ToLower(y)

	toID := map[string]OnMissingAction{
		"ignore": Ignore,
		"inject": Inject,
	}

	*a = toID[y]

	return nil
}

func (a OnMissingAction) MarshalYAML() (interface{}, error) {
	return a.String(), nil
}
