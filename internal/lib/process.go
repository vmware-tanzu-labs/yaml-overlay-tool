// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"

	"github.com/op/go-logging"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

var (
	log              = logging.MustGetLogger("lib") //nolint:gochecknoglobals
	ErrInvalidAction = errors.New("invalid overlay action")
)

func Process(options *Options) error {
	instructions, err := ReadInstructionFile(&options.InstructionsFile)
	if err != nil {
		return err
	}

	for fileIndex, file := range instructions.YamlFiles {
		for nodeIndex := range file.Nodes {
			log.Infof("Processing Common Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for i := range instructions.CommonOverlays {
				instructions.CommonOverlays[i].process(&instructions.YamlFiles[fileIndex], nodeIndex)
			}

			log.Infof("Processing File Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for i := range file.Overlays {
				file.Overlays[i].process(&instructions.YamlFiles[fileIndex], nodeIndex)
			}

			log.Infof("Processing Document Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for docIndex, doc := range file.Documents {
				if doc.Path != fmt.Sprint(docIndex) {
					continue
				}

				for i := range doc.Overlays {
					file.Documents[docIndex].Overlays[i].process(&instructions.YamlFiles[fileIndex], nodeIndex)
				}
			}

			output, err := yaml.Marshal(file.Nodes[nodeIndex])
			if err != nil {
				log.Errorf("unable to marshal final document %s, error: %s", file.Path, err)
			}

			log.Noticef("Final: >>>\n%s\n", output)
			// added so we can quickly see the results of the run
			fmt.Printf("---\n%s", output)
			// Save to File
		}
	}

	return nil
}

func (o *Overlay) process(f *YamlFile, i int) {
	if ok := checkIndex(i, o.DocumentIndex); !ok {
		return
	}

	node := f.Nodes[i]

	log.Debugf("%s at %s in file %s on Document %d\n", o.Action, o.Query, f.Path, i)

	yp, err := yamlpath.NewPath(o.Query)
	if err != nil {
		log.Errorf("an error occurred parsing the query path %v\n%w", o.Query, err)
	}

	results, err := yp.Find(node)
	if err != nil || results == nil {
		log.Debugf("Call OnMissing Here")
	}

	if err := processActions(node, results, o); err != nil {
		log.Errorf("%w in file %s on Document %d", err, f.Path, i)
	}
}

func checkIndex(current int, filter []int) bool {
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
				return fmt.Errorf("%w, skipping replace for %s result %d", err, o.Query, ri)
			}
		case "format":
			if err := actions.Format(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping format for %s result %d", err, o.Query, ri)
			}
		case "merge":
			if err := actions.Merge(results[ri], &o.Value); err != nil {
				return fmt.Errorf("%w, skipping merge for %s result %d", err, o.Query, ri)
			}
		default:
			return fmt.Errorf("%w: %s", ErrInvalidAction, o.Action)
		}
	}

	return nil
}
