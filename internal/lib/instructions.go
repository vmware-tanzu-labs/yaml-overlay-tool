// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func ReadInstructionFile(fileName *string) (*Instructions, error) {
	var instructions Instructions

	fmt.Printf("Instructions File: %s\n\n", *fileName)

	reader, err := ReadStream(*fileName)
	if err != nil {
		return nil, err
	}

	dc := yaml.NewDecoder(reader)
	if err := dc.Decode(&instructions); err != nil {
		return nil, err
	}

	fmt.Printf("%v\n\n", instructions)

	if err := instructions.ReadYamlFiles(); err != nil {
		return nil, err
	}

	return &instructions, nil
}
