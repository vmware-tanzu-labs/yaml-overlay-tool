package lib

import (
	"fmt"
	"log"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/edit"
	"gopkg.in/yaml.v3"
)

func Process(instructions *Instructions) error {
	for fileIndex, file := range instructions.YamlFiles {
		for nodeIndex := range file.Nodes {
			log.Printf("Processing Common Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for i := range instructions.CommonOverlays {
				instructions.CommonOverlays[i].process(&instructions.YamlFiles[fileIndex], nodeIndex)
			}

			log.Printf("Processing File Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for i := range file.Overlays {
				file.Overlays[i].process(&instructions.YamlFiles[fileIndex], nodeIndex)
			}

			log.Printf("Processing Document Overlays in File %s on Document %d\n\n", file.Path, nodeIndex)

			for docIndex, doc := range file.Documents {
				if doc.Path != fmt.Sprint(docIndex) {
					continue
				}

				for i := range doc.Overlays {
					file.Documents[docIndex].Overlays[i].process(&instructions.YamlFiles[fileIndex], nodeIndex)
				}
			}
		}
	}

	return nil
}

func (o *Overlay) process(f *YamlFile, i int) {
	var indexFound = true
	if o.DocumentIndex != nil {
		indexFound = false

		for di := range o.DocumentIndex {
			if i == o.DocumentIndex[di] {
				indexFound = true
				break
			}
		}
	}

	if !indexFound {
		return
	}

	var node = f.Nodes[i]

	log.Printf("%s at %s in file %s on Document %d\n", o.Action, o.Query, f.Path, i)

	result, err := edit.IteratePath(node, o.Query)
	if err != nil {
		log.Println("Call OnMissing Here")
	}

	b, _ := yaml.Marshal(&result)
	p, _ := yaml.Marshal(o.Value)

	log.Println("Current:")
	log.Println(string(b))
	log.Println("Proposed:")
	log.Println(string(p))
}
