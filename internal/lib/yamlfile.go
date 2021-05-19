// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlFile struct {
	Name      string      `yaml:"name,omitempty"`
	Overlays  []*Overlay  `yaml:"overlays,omitempty"`
	Documents []*Document `yaml:"documents,omitempty"`
	Files     Files       `yaml:"path,omitempty"`
}

func (yf *YamlFile) queueOverlays(oChan chan *workStream) {
	for _, f := range yf.Files {
		for nodeIndex, n := range f.Nodes {
			for _, o := range yf.Overlays {
				if ok := o.checkDocumentIndex(nodeIndex); ok {
					oChan <- &workStream{
						Overlay:   *o,
						Node:      n,
						NodeIndex: nodeIndex,
						Path:      f.Path,
					}
				}
			}

			for _, d := range yf.Documents {
				if ok := d.checkDocumentIndex(nodeIndex); ok {
					for _, o := range d.Overlays {
						oChan <- &workStream{
							Overlay:   *o,
							Node:      n,
							NodeIndex: nodeIndex,
							Path:      f.Path,
						}
					}
				}
			}
		}
	}

	close(oChan)
}

func (yf *YamlFile) doPostProcessing(cfg *Config) error {
	var o *os.File

	var err error

	output := new(bytes.Buffer)

	ye := yaml.NewEncoder(output)
	defer func() {
		if err = ye.Close(); err != nil {
			log.Fatalf("error closing encoder, %w", err)
		}
	}()

	ye.SetIndent(cfg.Indent)

	output.WriteString("---\n")

	for _, f := range yf.Files {
		for _, node := range f.Nodes {
			err = ye.Encode(node)
			if err != nil {
				return fmt.Errorf("unable to marshal final document %s, error: %w", f.Path, err)
			}
		}

		log.Debugf("Final: >>>\n%s\n", output)
		// added so we can quickly see the results of the run
		if cfg.StdOut {
			o = os.Stdout
		} else {
			o, err = f.OpenOutputFile(cfg)
			if err != nil {
				return err
			}
		}

		if _, err = output.WriteTo(o); err != nil {
			return fmt.Errorf("failed to %w", err)
		}

		output.Reset()
	}

	return nil
}
