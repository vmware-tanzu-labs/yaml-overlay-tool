// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type File struct {
	Nodes      []*yaml.Node
	Origin     string
	Path       string
	outputPath string
}

type Files []*File

func (f *Files) readYamlFile(p string) error {
	var files []string

	if ok, err := isDirectory(p); ok {
		files, err = getFileNames(p)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}

		files = []string{p}
	}

	for _, file := range files {
		source := &File{
			Origin: "file",
			Path:   file,
		}

		reader, err := ReadStream(file)
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
				return fmt.Errorf("failed to read file %s: %w", file, err)
			}

			source.Nodes = append(source.Nodes, &y)
		}

		*f = append(*f, source)
	}

	return nil
}

func (f *Files) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var p string

	if err := unmarshal(&p); err != nil {
		return err
	}

	return f.readYamlFile(p)
}

func (f *File) OpenFile(o *Config) (*os.File, error) {
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

	return file, nil
}
