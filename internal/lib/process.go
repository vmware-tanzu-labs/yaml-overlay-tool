package lib

import (
	"fmt"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/edit"
	"gopkg.in/yaml.v3"
)

func Process(instructions *Instructions) error {
	for fileIndex, file := range instructions.YamlFiles {
		for i := range instructions.CommonOverlays {
			instructions.CommonOverlays[i].process(&instructions.YamlFiles[fileIndex])
		}

		for i := range file.Overlays {
			file.Overlays[i].process(&instructions.YamlFiles[fileIndex])
		}

		// for docIndex, doc := range file.Documents {
		//	for i := range doc.Overlays {
		//		file.Documents[docIndex].Overlays[i].process(&instructions.YamlFiles[fileIndex].Documents[docIndex])
		//	}
		//}
	}

	return nil
}

func (o *Overlay) process(f *YamlFile) {
	var node = f.Node

	fmt.Printf("%s at %s on file %s\n", o.Action, o.Query, f.Path)

	result, err := edit.IteratePath(node, o.Query)
	if err != nil {
		fmt.Println("Call OnMissing Here")
	}

	b, _ := yaml.Marshal(&result)
	fmt.Println(string(b))
}
