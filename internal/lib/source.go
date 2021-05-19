// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Source struct {
	Nodes      []*yaml.Node
	Origin     string
	Path       string
	outputPath string
}

type Sources []*Source

func (src *Sources) readYamlFile(p string) error {
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
		source := &Source{
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

		*src = append(*src, source)
	}

	return nil
}

func (src *Sources) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var p string

	if err := unmarshal(&p); err != nil {
		return err
	}

	return src.readYamlFile(p)
}

func (src *Source) Save(o *Config, output string) error {
	fileName := path.Join(o.OutputDir, "yamlFiles", src.outputPath)
	dirName := path.Dir(fileName)

	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.MkdirAll(dirName, 0755); err != nil {
			return fmt.Errorf("failed to create output directory %s, %w", dirName, err)
		}
	}

	//nolint:gosec //output files with read and write permissions so that end-users can continue to leverage these files
	if err := os.WriteFile(fileName, []byte(output), 0644); err != nil {
		return fmt.Errorf("failed to save file %s: %w", fileName, err)
	}

	return nil
}

func (src *Source) doPostProcessing(cfg *Config) error {
	output := new(bytes.Buffer)
	ye := yaml.NewEncoder(output)
	ye.SetIndent(cfg.Indent)

	for _, node := range src.Nodes {
		err := ye.Encode(node)
		if err != nil {
			return fmt.Errorf("unable to marshal final document %s, error: %w", src.Path, err)
		}
	}

	final := fmt.Sprintf("---\n%s\n", output)

	log.Noticef("Final: >>>\n%s\n", output)
	// added so we can quickly see the results of the run
	if cfg.StdOut {
		writer := bufio.NewWriter(os.Stdout)
		if _, err := writer.WriteString(final); err != nil {
			return fmt.Errorf("failed to write to string, %w", err)
		}

		writer.Flush()
	}

	if err := src.Save(cfg, final); err != nil {
		return fmt.Errorf("failed to save, %w", err)
	}

	return nil
}
