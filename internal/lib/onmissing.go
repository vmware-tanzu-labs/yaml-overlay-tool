// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/builder"
	"gopkg.in/yaml.v3"
)

var (
	ErrOnMissingNoInjectAction = errors.New("no matches and no onMissing.action of 'inject'")
	ErrOnMissingNoInjectPath   = errors.New("no matches and no onMissing.injectPath")
	ErrOnMissingInvalidType    = errors.New("invalid type for onMissing.injectPath")
)

type OnMissing struct {
	Action     OnMissingAction `yaml:"action,omitempty"`
	InjectPath multiString     `yaml:"injectPath,omitempty"`
}

func (o *Overlay) onMissing(n *yaml.Node) error {
	// check if the query has a match
	// if no match then we require an inject path
	// we need to then check if each inject path is valid (does it exist)
	// if we had an inject path(s) then we inject the value to those locations
	// if we didn't have an inject path we have an implicit onMissing: ignore and we put out a warning if not stdout option to terminal
	switch o.OnMissing.Action {
	case Inject:
		return o.handleInjectPath(n)
	default:
		log.Debugf("ignoring %s at %s due to %s\n", o.Action, o.Query, ErrOnMissingNoInjectAction)

		return nil
	}
}

func (o *Overlay) doInjectPath(ip []string, node *yaml.Node) error {
	bps, err := builder.NewPaths(ip)
	if err != nil {
		return fmt.Errorf("failed to build inject path %s, %w", ip, err)
	}

	y, err := bps.BuildPaths()
	if err != nil {
		return fmt.Errorf("failed to build inject path %s, %w", ip, err)
	}

	err = actions.Merge(node, y)
	if err != nil {
		return fmt.Errorf("failed to merge injectpath scaffolding %s with document, %w", ip, err)
	}

	results, err := searchYAMLPaths(ip, node)
	if err != nil {
		return fmt.Errorf("%w, on injectPath %s", err, ip)
	}

	for _, r := range results {
		if err := actions.Replace(r, &o.Value); err != nil {
			return fmt.Errorf("%w for onMissing.InjectPath", err)
		}
	}

	return nil
}

func (o *Overlay) handleInjectPath(n *yaml.Node) error {
	_, err := builder.NewPaths(o.Query)
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
