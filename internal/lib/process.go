// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"

	"github.com/op/go-logging"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

var (
	log              = logging.MustGetLogger("lib") //nolint:gochecknoglobals
	ErrInvalidAction = errors.New("invalid overlay action")
)

func Execute(options *Options) error {
	instructions, err := ReadInstructionFile(&options.InstructionsFile)
	if err != nil {
		return err
	}

	return applyOverlays(instructions, options)
}

func checkDocIndex(current int, filter []int) bool {
	if filter != nil {
		for f := range filter {
			if current == filter[f] {
				return true
			}
		}

		return false
	}

	return true
}

func processActions(node *yaml.Node, results []*yaml.Node, o *Overlay) error {
	for ri := range results {
		b, _ := yaml.Marshal(&results[ri])
		p, _ := yaml.Marshal(o.Value)

		log.Debugf("Current: >>>\n%s\n", b)
		log.Debugf("Proposed: >>>\n%s\n", p)

		// do something with the results based on the provided overlay action
		switch o.Action {
		case "delete":
			actions.Delete(node, results[ri])
		case "replace":
			if err := actions.Replace(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping replace for %s query result %d", err, o.Query, ri)
			}
		case "format":
			if err := actions.Format(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping format for %s query result %d", err, o.Query, ri)
			}
		case "merge":
			if err := actions.Merge(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping merge for %s query result %d", err, o.Query, ri)
			}
		default:
			return fmt.Errorf("%w of type '%s'", ErrInvalidAction, o.Action)
		}
	}

	return nil
}

func applyOverlays(instructions *Instructions, options *Options) error {
	for _, file := range instructions.YamlFiles {
		for nodeIndex := range file.Nodes {
			log.Infof("Processing Common Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			if err := file.processOverlays(instructions.CommonOverlays, nodeIndex); err != nil {
				return fmt.Errorf("failed to apply common overlays, %w", err)
			}

			log.Infof("Processing File Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			if err := file.processOverlays(file.Overlays, nodeIndex); err != nil {
				return fmt.Errorf("failed to apply file overlays, %w", err)
			}

			log.Infof("Processing Document Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for docIndex, doc := range file.Documents {
				if doc.Path != fmt.Sprint(docIndex) {
					continue
				}

				if err := file.processOverlays(file.Documents[docIndex].Overlays, nodeIndex); err != nil {
					return fmt.Errorf("failed to apply document overlays, %w", err)
				}
			}
		}

		if err := file.doPostProcessing(options); err != nil {
			return fmt.Errorf("failed to preform post processing on %s: %w", file.Path, err)
		}
	}

	return nil
}
