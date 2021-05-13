// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func (yf *YamlFile) readYamlFile(path string) error {
	source := &Source{
		Origin: "file",
		Path:   path,
	}

	reader, err := ReadStream(path)
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

		source.Nodes = append(source.Nodes, &y)
	}

	yf.Source = append(yf.Source, *source)

	return nil
}
