// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/builder"
	"gopkg.in/yaml.v3"
)

var (
	// ErrOnMissingNoInjectAction occurrs if there are not matches found and no inject action present.
	ErrOnMissingNoInjectAction = errors.New("no matches and no onMissing.action of 'inject'")

	// ErrOnMissingNoInjectPath occurrs if there are no matches found for the injectPath.
	ErrOnMissingNoInjectPath = errors.New("no matches and no onMissing.injectPath")
)

type OnMissing struct {
	Action     actions.OnMissingAction `yaml:"action,omitempty"`
	InjectPath Queries                 `yaml:"injectPath,omitempty"`
}

func (o *Overlay) onMissing(n *yaml.Node) error {
	if o.OnMissing.Action == actions.Default {
		if strings.EqualFold(viper.GetString("defaultOnMissingAction"), "inject") {
			o.OnMissing.Action = actions.Inject
		}
	}

	// check if the query has a match
	// if no match then we require an inject path
	// we need to then check if each inject path is valid (does it exist)
	// if we had an inject path(s) then we inject the value to those locations
	// if we didn't have an inject path we have an implicit onMissing: ignore and we put out a warning if not stdout option to terminal
	switch o.OnMissing.Action {
	case actions.Inject:
		return o.handleInjectPath(n)
	default:
		log.Debugf("ignoring %s at %s due to %s\n", o.Action, o.Query, ErrOnMissingNoInjectAction)

		return nil
	}
}

func (o *Overlay) doInjectPath(ip Queries, node *yaml.Node) error {
	bps, err := builder.NewPaths(ip.Paths())
	if err != nil {
		return fmt.Errorf("failed to build inject path %s, %w", ip, err)
	}

	y, _ := bps.BuildPaths()

	err = actions.MergeNode(node, y)
	if err != nil {
		return fmt.Errorf("failed to merge injectpath scaffolding %s with document, %w", ip, err)
	}

	results := ip.Find(node)

	for _, r := range results {
		if err := actions.ReplaceNode(r, &o.Value); err != nil {
			return fmt.Errorf("%w for onMissing.InjectPath", err)
		}
	}

	return nil
}

func (o *Overlay) handleInjectPath(n *yaml.Node) error {
	_, err := builder.NewPaths(o.Query.Paths())
	if err != nil {
		if errors.Is(err, builder.ErrInvalidPathSyntax) {
			if o.OnMissing.InjectPath == nil {
				log.Debugf("ignoring %s at %s due to %s\n", o.Action, o.Query, ErrOnMissingNoInjectPath)

				return nil
			}

			return o.doInjectPath(o.OnMissing.InjectPath, n)
		}

		return fmt.Errorf("%w, for onMissing", err)
	}

	return o.doInjectPath(o.Query, n)
}
