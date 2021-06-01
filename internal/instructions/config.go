// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/op/go-logging"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

// Config contains configuration options used with instruction files.
type Config struct {
	Verbose          bool
	LogLevel         logging.Level
	InstructionsFile string
	OutputDir        string
	StdOut           bool
	Indent           int
	Styles           actions.Styles
	Values           []string
}

// doPostProcessing renders a document and outputs it to the location specified in config.
func (cfg *Config) doPostProcessing(yf *YamlFile) error {
	var o *os.File

	var err error

	output, err := cfg.encodeNodes(yf.Nodes)
	if err != nil {
		return fmt.Errorf("unable to marshal final document %s, error: %w", yf.Path, err)
	}

	if output.Len() == 0 {
		log.Debugf("File %s was omitted from output", yf.Path)

		return nil
	}

	if cfg.StdOut {
		o = os.Stdout
	} else {
		log.Debugf("Final: >>>\n%s\n", output)
		o, err = cfg.openOutputFile(yf)
		if err != nil {
			return err
		}

		defer CloseFile(o)
	}

	if _, err = output.WriteTo(o); err != nil {
		return fmt.Errorf("failed to %w", err)
	}

	output.Reset()

	return nil
}

// openOutputFile opens or creates a file for outputing results.
func (cfg *Config) openOutputFile(yf *YamlFile) (*os.File, error) {
	fileName := path.Join(cfg.OutputDir, yf.OutputPath)
	dirName := path.Dir(fileName)

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return nil, fmt.Errorf("failed to create output directory %s, %w", dirName, err)
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create/open file %s: %w", fileName, err)
	}

	os.Stdout.WriteString(fileName + "\n")

	return file, nil
}

// encodeNodes processes each node given and encodes them into a single buffer for output.
func (cfg *Config) encodeNodes(nodes []*yaml.Node) (*bytes.Buffer, error) {
	var fileWritten bool

	output := new(bytes.Buffer)

	ye := yaml.NewEncoder(output)

	defer func() {
		if fileWritten {
			if err := ye.Close(); err != nil {
				log.Criticalf("error closing encoder, %s", err)
			}
		}
	}()

	ye.SetIndent(cfg.Indent)

	for i, node := range nodes {
		if len(node.Content) == 0 {
			continue
		}

		if i == 0 {
			output.WriteString("---\n")
		}

		actions.SetStyle(cfg.Styles, node)

		err := ye.Encode(node)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		fileWritten = true
	}

	return output, nil
}
