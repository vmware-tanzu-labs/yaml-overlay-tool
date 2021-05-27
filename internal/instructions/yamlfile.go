// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/overlays"
	"gopkg.in/yaml.v3"
)

var ErrFoundDirectoryWithPathOutput = errors.New("found directory as input path with file as output path")

var ErrAbsolutePathForOutputPath = errors.New("absolute paths are currently not supported for outputPath")

// YamlFile is used to define which files should be manipulated and overlays specific to that file.
type YamlFile struct {
	// Optional Name to define for organization purposes.
	Name string `yaml:"name,omitempty"`
	// List of Overlays specific to this yamlFile.
	Overlays []*overlays.Overlay `yaml:"overlays,omitempty"`
	// a list of more specific entries and overlays for a specific document within the yamlFile.
	Documents []*Document `yaml:"documents,omitempty"`

	Nodes []*yaml.Node

	Path string `yaml:"path,omitempty"`

	OutputPath string `yaml:"outputPath,omitempty"`
}

type YamlFiles []*YamlFile

// queueOverlays sends all overlay jobs to the workstream for processing.
func (yf *YamlFile) queueOverlays(stream *overlays.WorkStream) {
	for nodeIndex, n := range yf.Nodes {
		for _, o := range yf.Overlays {
			if ok := o.CheckDocumentIndex(nodeIndex); ok {
				stream.AddWorkload(o, n, nodeIndex, yf.Path)
			}
		}

		for _, d := range yf.Documents {
			if ok := d.checkDocumentIndex(nodeIndex); ok {
				for _, o := range d.Overlays {
					stream.AddWorkload(o, n, nodeIndex, yf.Path)
				}
			}
		}
	}

	stream.CloseStream()
}

// doPostProcessing renders a document and outputs it to the location specified in config.
func (yf *YamlFile) doPostProcessing(cfg *Config) error {
	var o *os.File

	var err error

	var fileWritten bool

	output := new(bytes.Buffer)

	ye := yaml.NewEncoder(output)

	defer func() {
		if fileWritten {
			if err = ye.Close(); err != nil {
				log.Criticalf("error closing encoder, %s", err)
			}
		}
	}()

	ye.SetIndent(cfg.Indent)

	for i, node := range yf.Nodes {
		if len(node.Content) == 0 {
			continue
		}

		if i == 0 {
			output.WriteString("---\n")
		}

		actions.SetStyle(cfg.Styles, node)

		err = ye.Encode(node)
		if err != nil {
			return fmt.Errorf("unable to marshal final document %s, error: %w", yf.Path, err)
		}

		fileWritten = true
	}

	// added so we can quickly see the results of the run
	if cfg.StdOut {
		o = os.Stdout
	} else {
		log.Debugf("Final: >>>\n%s\n", output)
		o, err = yf.OpenOutputFile(cfg)
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

func (yfs *YamlFiles) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var yft []*YamlFile

	if err := unmarshal(&yft); err != nil {
		return err
	}

	y := YamlFiles(yft)
	if err := y.expandDirectories(); err != nil {
		return err
	}

	for _, yf := range y {
		if err := yf.readYamlFile(); err != nil {
			return err
		}
	}

	*yfs = y

	return nil
}

// OpenOutputFile opens or creates a file for outputing results.
func (yf *YamlFile) OpenOutputFile(o *Config) (*os.File, error) {
	fileName := path.Join(o.OutputDir, yf.OutputPath)
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

func (yf *YamlFile) readYamlFile() error {
	reader, err := ReadStream(yf.Path)
	if err != nil {
		return err
	}

	dc := yaml.NewDecoder(reader)

	for {
		var y yaml.Node

		if err := dc.Decode(&y); errors.Is(err, io.EOF) {
			if reader, ok := reader.(*os.File); ok {
				CloseFile(reader)

				break
			}
		} else if err != nil {
			return fmt.Errorf("failed to read file %s: %w", yf.Path, err)
		}

		yf.Nodes = append(yf.Nodes, &y)
	}

	return nil
}

func (yfs *YamlFiles) expandDirectories() error {
	y := []*YamlFile(*yfs)

	var paths []string

	var removeItems []int

	for i, yf := range y {
		if !path.IsAbs(yf.Path) {
			yf.Path = path.Join(instructionsDir, yf.Path)
		}

		if path.IsAbs(yf.OutputPath) {
			return ErrAbsolutePathForOutputPath
		}

		if ok, err := isDirectory(yf.Path); err != nil {
			return err
		} else if ok {
			if path.Ext(yf.OutputPath) != "" {
				return fmt.Errorf("%w, provide a directory for outputPath instead", ErrFoundDirectoryWithPathOutput)
			}
			paths, err = getFileNames(yf.Path)
			if err != nil {
				return err
			}

			removeItems = append(removeItems, i)

			for _, p := range paths {
				sp := &YamlFile{
					Name:       yf.Name,
					Overlays:   yf.Overlays,
					Documents:  yf.Documents,
					Path:       p,
					OutputPath: yf.OutputPath,
				}

				y = append(y, sp)
			}
		}
	}

	for _, remove := range removeItems {
		y[remove] = y[len(y)-1]
		y = y[:len(y)-1]
	}

	*yfs = y

	return nil
}
