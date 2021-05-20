// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"bytes"
	"fmt"
	"os"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/overlays"
	"gopkg.in/yaml.v3"
)

type YamlFile struct {
	Name      string              `yaml:"name,omitempty"`
	Overlays  []*overlays.Overlay `yaml:"overlays,omitempty"`
	Documents []*Document         `yaml:"documents,omitempty"`
	Files     Files               `yaml:"path,omitempty"`
}

func (yf *YamlFile) queueOverlays(stream *overlays.WorkStream) {
	for _, f := range yf.Files {
		for nodeIndex, n := range f.Nodes {
			for _, o := range yf.Overlays {
				if ok := o.CheckDocumentIndex(nodeIndex); ok {
					stream.AddWorkload(o, n, nodeIndex, f.Path)
				}
			}

			for _, d := range yf.Documents {
				if ok := d.checkDocumentIndex(nodeIndex); ok {
					for _, o := range d.Overlays {
						stream.AddWorkload(o, n, nodeIndex, f.Path)
					}
				}
			}
		}
	}

	stream.CloseStream()
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

		// added so we can quickly see the results of the run
		if cfg.StdOut {
			o = os.Stdout
		} else {
			log.Debugf("Final: >>>\n%s\n", output)
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
