// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"fmt"

	"github.com/op/go-logging"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

var log = logging.MustGetLogger("lib")

func Process(instructions *Instructions) error {
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
		}
	}

	return nil
}

func (o *Overlay) process(f *YamlFile, i int) {
	var indexFound = true
	if o.DocumentIndex != nil {
		indexFound = false

		for di := range o.DocumentIndex {
			if i == o.DocumentIndex[di] {
				indexFound = true
				break
			}
		}
	}

	if !indexFound {
		return
	}

	var node = f.Nodes[i]

	log.Debugf("%s at %s in file %s on Document %d\n", o.Action, o.Query, f.Path, i)

	yp, err := yamlpath.NewPath(o.Query)
	if err != nil {
		log.Errorf("an error occurred parsing the query path %v\n%v", o.Query, err)
	}

	results, err := yp.Find(node)
	if err != nil || results == nil {
		log.Debugf("Call OnMissing Here")
	}

	for i := range results {
		b, _ := yaml.Marshal(&results[i])
		p, _ := yaml.Marshal(o.Value)

		log.Debugf("Current: >>>\n%s\n", b)
		log.Debugf("Proposed: >>>\n%s\n", p)

		// do something with the results based on the provided overlay action
		switch o.Action {
		case "delete":
			actions.Delete(node, results[i])
		case "replace":
			actions.Replace(results[i], &o.Value)
		case "merge":
		default:
			log.Errorf("Invalid overlay action: %v", o.Action)
		}
	}
}
