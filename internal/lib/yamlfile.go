package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v3"
)

func (f *YamlFile) processOverlays(o []Overlay, nodeIndex int) error {
	for i := range o {
		if err := o[i].process(f, nodeIndex); err != nil {
			return err
		}
	}

	return nil
}

func (f *YamlFile) Save(o *Options) error {
	b, err := yaml.Marshal(f.Nodes)
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %w", f.Path, err)
	}

	fileName := path.Base(f.Path)
	outputFileName := path.Join(o.OutputDir, fileName)
	//nolint:gosec //output files with read and write permissions so that end-users can continue to leverage these files
	if err := ioutil.WriteFile(outputFileName, b, 0644); err != nil {
		return fmt.Errorf("failed to save file %s: %w", outputFileName, err)
	}

	return nil
}

func (f *YamlFile) doPostProcessing(o *Options) error {
	output := new(bytes.Buffer)
	ye := yaml.NewEncoder(output)
	ye.SetIndent(o.Indent)

	for _, node := range f.Nodes {
		err := ye.Encode(node)
		if err != nil {
			return fmt.Errorf("unable to marshal final document %s, error: %w", f.Path, err)
		}
	}

	log.Noticef("Final: >>>\n%s\n", output)
	// added so we can quickly see the results of the run
	if o.StdOut {
		fmt.Printf("---\n%s", output) //nolint:forbidigo
	}

	return nil
}
