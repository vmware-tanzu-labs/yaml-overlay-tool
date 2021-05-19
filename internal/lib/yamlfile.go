// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

type Document struct {
	Name     string     `yaml:"name,omitempty"`
	Path     int        `yaml:"index,omitempty"`
	Overlays []*Overlay `yaml:"overlays,omitempty"`
}

type YamlFile struct {
	Name      string      `yaml:"name,omitempty"`
	Overlays  []*Overlay  `yaml:"overlays,omitempty"`
	Documents []*Document `yaml:"documents,omitempty"`
	Sources   Sources     `yaml:"path,omitempty"`
}

func (d *Document) checkDocumentIndex(docIndex int) bool {
	return d.Path == docIndex
}

func (yf *YamlFile) queueSourceFiles(oChan chan *workStream) {
	for _, src := range yf.Sources {
		for nodeIndex := range src.Nodes {
			for _, o := range yf.Overlays {
				if ok := o.checkDocumentIndex(nodeIndex); ok {
					oChan <- &workStream{
						Overlay:   *o,
						NodeIndex: nodeIndex,
						File:      src,
					}
				}
			}

			for _, d := range yf.Documents {
				if ok := d.checkDocumentIndex(nodeIndex); ok {
					for _, o := range d.Overlays {
						oChan <- &workStream{
							Overlay:   *o,
							NodeIndex: nodeIndex,
							File:      src,
						}
					}
				}
			}
		}
	}

	close(oChan)
}

func (i *Instructions) queueYamlFiles(c chan *YamlFile) {
	for _, yf := range i.YamlFiles {
		c <- yf
	}

	close(c)
}

func OverlayHandler(cfg *Config, oChan chan *workStream, errs chan error) {
	for o := range oChan {
		if err := o.Overlay.applyOverlay(o.File, o.NodeIndex); err != nil {
			errs <- err
		}
	}
}
