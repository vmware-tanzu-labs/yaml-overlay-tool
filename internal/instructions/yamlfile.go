// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/viper"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/overlays"
	"gopkg.in/yaml.v3"
)

var ErrFoundDirectoryWithPathOutput = errors.New("found directory as input path with file as output path")

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

func (yfs *YamlFiles) UnmarshalYAML(value *yaml.Node) error {
	var yft []*YamlFile

	if err := value.Decode(&yft); err != nil {
		return fmt.Errorf("%w at line %d column %d", err, value.Line, value.Column)
	}

	y := YamlFiles(yft)
	if err := y.expandDirectories(); err != nil {
		return fmt.Errorf("%w at line %d column %d", err, value.Line, value.Column)
	}

	y.mergeDuplicates()

	for _, yf := range y {
		if err := yf.readYamlFile(); err != nil {
			return err
		}
	}

	*yfs = y

	return nil
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
			}

			break
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

	for i := 0; i <= len(y)-1; i++ {
		if y[i].Path == "-" {
			continue
		}

		if !path.IsAbs(y[i].Path) {
			y[i].Path = path.Join(viper.GetString("instructionsDir"), y[i].Path)
		}

		if ok, err := isDirectory(y[i].Path); err != nil {
			return err
		} else if ok {
			if path.Ext(y[i].OutputPath) != "" {
				return fmt.Errorf("%w, provide a directory for outputPath instead", ErrFoundDirectoryWithPathOutput)
			}
			paths, err = getFileNames(y[i].Path)
			if err != nil {
				return err
			}

			for _, p := range paths {
				sp := &YamlFile{
					Name:       y[i].Name,
					Overlays:   y[i].Overlays,
					Documents:  y[i].Documents,
					Path:       p,
					OutputPath: y[i].OutputPath,
				}

				y = append(y, sp)
			}

			y = append(y[:i], y[i+1:]...)
			i--
		}
	}

	*yfs = y

	return nil
}

func (yfs *YamlFiles) mergeDuplicates() {
	ys := []*YamlFile(*yfs)

	search := make(map[string]*YamlFile, len(*yfs))

	for i := 0; i <= len(ys)-1; i++ {
		if search[ys[i].Path] == nil {
			search[ys[i].Path] = ys[i]
		} else {
			search[ys[i].Path].Overlays = append(search[ys[i].Path].Overlays, ys[i].Overlays...)
			search[ys[i].Path].Documents = append(search[ys[i].Path].Documents, ys[i].Documents...)
			if ys[i].OutputPath != "" {
				search[ys[i].Path].OutputPath = ys[i].OutputPath
			}

			ys = append(ys[:i], ys[i+1:]...)
			i--
		}
	}

	*yfs = ys
}

// removeCommentsFromNode will traverse the tree of yaml nodes and remove all comments.
func removeCommentsFromNode(node *yaml.Node) {
	if node.Content == nil {
		return
	}

	for _, child := range node.Content {
		child.HeadComment = ""
		child.LineComment = ""
		child.FootComment = ""

		removeCommentsFromNode(child)
	}
}
