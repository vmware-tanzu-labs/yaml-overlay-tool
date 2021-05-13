package lib

import (
	"errors"
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
)

func (a Action) String() string {
	toString := map[Action]string{
		Invalid: "",
		Merge:   "merge",
		Replace: "replace",
		Delete:  "delete",
		Format:  "format",
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
		"format":  Format,
	}

	if toID[y] == Invalid {
		return ErrInvalidAction
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

	if toID[y] == Invalid {
		return ErrInvalidAction
	}

	*a = toID[y]

	return nil
}

func (a OnMissingAction) MarshalYAML() (interface{}, error) {
	return a.String(), nil
}
