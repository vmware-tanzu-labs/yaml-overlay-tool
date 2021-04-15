package lib

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func ReadInstructionFile(fileName *string) error {
	var instructions Instructions

	fmt.Printf("Instructions File: %s\n\n", *fileName)

	reader, err := ReadStream(*fileName)
	if err != nil {
		return err
	}

	dc := yaml.NewDecoder(reader)
	if err := dc.Decode(&instructions); err != nil {
		return err
	}

	fmt.Println(instructions)

	if err := instructions.ReadYamlFiles(); err != nil {
		return err
	}

	fmt.Print(instructions.YamlFiles[0].Node)

	return nil
}
