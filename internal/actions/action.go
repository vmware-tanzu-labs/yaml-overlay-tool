// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidAction occurs when user passes a action that is not one of merge, replace, delete, combine.
var ErrInvalidAction = errors.New("invalid overlay action")

type Action int

const (
	Invalid = iota
	// Merge overlay action.
	Merge
	// Replace overlay action.
	Replace
	// Delete overlay action.
	Delete
	// Combine overlay action.
	Combine
)

func (a Action) String() string {
	toString := map[Action]string{
		Invalid: "",
		Merge:   "merge",
		Replace: "replace",
		Delete:  "delete",
		Combine: "combine",
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
		"combine": Combine,
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
	// Ignore onMissing action.
	Ignore = iota
	// Inject onMissing action.
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
