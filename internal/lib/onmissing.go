// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/path"
	"gopkg.in/yaml.v3"
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
	case Ignore:
		log.Debugf("ignoring %s at %s due to %s\n", o.Action, o.Query, ErrOnMissingNoInjectAction)

		return nil
	case Inject:
		_, err := path.BuildMulti(o.Query)
		if err != nil {
			if errors.Is(err, path.ErrInvalidPathSyntax) {
				return o.handleInjectPath(n)
			}

			return fmt.Errorf("%w, for onMissing", err)
		}

		return o.doInjectPath(o.Query, n)
	default:
		return fmt.Errorf("%w for onMissing of type '%s'", ErrInvalidAction, o.Action)
	}
}

func (o *Overlay) doInjectPath(ip []string, node *yaml.Node) error {
	y, err := path.BuildMulti(ip)
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
			if errors.Is(err, ErrInvalidAction) {
				return fmt.Errorf("%w in instructions file", err)
			}

			return fmt.Errorf("%w for onMissing.InjectPath", err)
		}
	}

	return nil
}

func (o *Overlay) handleInjectPath(n *yaml.Node) error {
	_, err := path.BuildMulti(o.Query)
	if !errors.Is(err, path.ErrInvalidPathSyntax) {
		return o.doInjectPath(o.Query, n)
	}

	if o.OnMissing.InjectPath == nil {
		log.Debugf("ignoring %s at %s due to %s\n", o.Action, o.Query, ErrOnMissingNoInjectPath)

		return nil
	}

	return o.doInjectPath(o.OnMissing.InjectPath, n)
}
