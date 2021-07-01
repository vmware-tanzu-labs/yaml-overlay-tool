// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/op/go-logging"
	"gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("overlays") //nolint:gochecknoglobals

// ErrInvalidAction occurs when user passes a action that is not one of merge, replace, delete, combine.
var ErrInvalidAction = errors.New("invalid overlay action")

type Action int

const (
	// Invalid overlay action.
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

func (a *Action) UnmarshalYAML(value *yaml.Node) error {
	var y string

	if err := value.Decode(&y); err != nil {
		return fmt.Errorf("%w at line %d column %d", err, value.Line, value.Column)
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

func (a *Action) Set(val string) error {
	if err := yaml.Unmarshal([]byte(val), a); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (a *Action) Type() string {
	return "actions.Action"
}

type OnMissingAction int

const (
	// Default onMissing action. settable via viepr config, defaults to ignore.
	Default = iota
	// Ignore onMissing action.
	Ignore
	// Inject onMissing action.
	Inject
)

func (a OnMissingAction) String() string {
	toString := map[OnMissingAction]string{
		Default: "default",
		Ignore:  "ignore",
		Inject:  "inject",
	}

	return toString[a]
}

func (a *OnMissingAction) UnmarshalYAML(value *yaml.Node) error {
	var y string

	if err := value.Decode(&y); err != nil {
		return fmt.Errorf("%w at line %d column %d", err, value.Line, value.Column)
	}

	y = strings.ToLower(y)

	toID := map[string]OnMissingAction{
		"default": Default,
		"ignore":  Ignore,
		"inject":  Inject,
	}

	*a = toID[y]

	return nil
}

func (a OnMissingAction) MarshalYAML() (interface{}, error) {
	return a.String(), nil
}

func (a *OnMissingAction) Set(val string) error {
	if err := yaml.Unmarshal([]byte(val), a); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (a *OnMissingAction) Type() string {
	return "actions.OnMissingAction"
}
