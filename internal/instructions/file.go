// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

// File is unmarshalled yaml file and its path that overlay's will be applied to.
type File struct {
	Nodes      []*yaml.Node
	Origin     string
	Path       string
	outputPath string
}

// Files is a collection of file objects.
type Files []*File

// readYamlFile will resolve a list of file names from a path and read them into a files collection.
func (f *Files) readYamlFile(p string) error {
	var paths []string

	if ok, err := isDirectory(p); ok {
		paths, err = getFileNames(p)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}

		paths = []string{p}
	}

	for _, p := range paths {
		file := &File{
			Origin: "file",
			Path:   p,
		}

		reader, err := ReadStream(p)
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
				return fmt.Errorf("failed to read file %s: %w", p, err)
			}

			file.Nodes = append(file.Nodes, &y)
		}

		*f = append(*f, file)
	}

	return nil
}

// UnmarshalYAML is a custom unmarshal function to unmarshal Files objects.
func (f *Files) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var p string

	if err := unmarshal(&p); err != nil {
		return err
	}

	return f.readYamlFile(p)
}

// OpenOutputFile opens or creates a file for outputing results.
func (f *File) OpenOutputFile(o *Config) (*os.File, error) {
	fileName := path.Join(o.OutputDir, "yamlFiles", f.outputPath)
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
