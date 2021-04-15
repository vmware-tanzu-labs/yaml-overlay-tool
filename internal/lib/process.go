package lib

import "fmt"

func Process(instructions *Instructions) error {
	for fileIndex, file := range instructions.YamlFiles {
		for i := range instructions.CommonOverlays {
			instructions.CommonOverlays[i].process(&instructions.YamlFiles[fileIndex])
		}

		for i := range file.Overlays {
			file.Overlays[i].process(&instructions.YamlFiles[fileIndex])
		}

		for docIndex, doc := range file.Documents {
			for i := range doc.Overlays {
				file.Documents[docIndex].Overlays[i].process(&instructions.YamlFiles[fileIndex].Documents[docIndex])
			}
		}
	}

	return nil
}

func (o *Overlay) process(f *YamlFile) {
	fmt.Printf("%s at %s on file %s\n", o.Action, o.Query, f.Path)
}
