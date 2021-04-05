package commands

import (
	"fmt"
	"log"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/models"
	"gopkg.in/yaml.v2"
)

type StaticYamlFiles struct {
	Yaml_files []Name `yaml:"yaml_files,omitempty"`
}

type Name struct {
	Name     string     `yaml:"name"`
	Path     string     `yaml:"path"`
	Overlays []Overlays `yaml:"overlays"`
}

type Overlays struct {
	Name   string `yaml:"name"`
	Query  string `yaml:"query"`
	Value  string `yaml:"value"`
	Action string `yaml:"action"`
}

func instructFile(fn string) error {

	// teams full example to play with

	var data = `
---
# the commonOverlays apply to all yamlFiles listed out in 'yamlFiles' and are processed first on each file
commonOverlays:
  - name: "add a label to certain yaml documents with refined criteria"
    query: metadata.labels
    value: {'namespace': 'tanzu-dns'}
    action: merge
    # qualifier to further refine when this overlay is applied
    documentQuery:
      # default operator behavior is 'and' and has been omitted as an example of
      ## this behavior
      # all of the 'and' operator queries must match or any one of the 'or'
      ## operator queries
    - conditions:
      - key: kind
        value: Service
      - key: metadata.labels.'app.kubernetes.io/name'
        value: external-dns
    - conditions:
      - key: metadata.name
        value: pvc-var-cache-bind
  - name: "add a common label to everything"
    query: metadata.labels
    value: {'cool_label': 'bro'}
    action: merge
yamlFiles: # what to overlay onto
  - name: "some arbitrary descriptor" # Name is Optional
    path: "examples/manifests/test.yaml"
    overlays: # if multi-doc yaml file, applies to all docs, gets applied first
    - name: "delete all annotations"
      query: metadata.annotations
      value: {}
      action: "delete"
    - name: "add in a new label"
      query: metadata.labels
      value: {'some': 'thing'}
      action: "merge"
      onMissing:
        action: "inject" # inject | ignore
    - name: "Change the apiVersion to v2alpha1"
      query: apiVersion
      value: v2alpha1
      action: replace
    # on the following 2 items, notice that the onMissing is not set
    ## these will only affect the yaml docs that have matches, otherwise ignore
    - name: "Merge in a list item"
      query: spec.ports
      value:
        - name: dns-tcp
          port: 53
          protocol: TCP
          targetPort: dns-tcp
      action: merge
    # not really a real-world example, but showing off functionality
    - name: "now replace the merged list with just the new port"
      query: spec.ports
      value:
        - name: dns-tcp
          port: 53
          protocol: TCP
          targetPort: dns-tcp
      action: replace
    # next one shouldn't do anything because no onMissing = implicit ignore
    - query: status
      value: {}
      action: "merge"
    - name: "Demo the need for an inject path"
      query: fake.key1.*
      value: {'fake': 'content1'}
      action: "merge"
      onMissing:
        action: "inject"
    # same as previous, but with an injectPath (actually does this one)
    - name: "Show same example but with an injectPath"
      query: fake.key2.*
      value: {'fake': 'content2'}
      action: "merge"
      onMissing:
        action: "inject"
        injectPath: fake2.key2
      # qualifier to only apply to the first doc in the yaml file
      documentIndex:
        - 0
    documents: # optional and only used for multi-doc yaml files
    # need to refer to them by their index
    - name: the manifest that does something
      path: 0
      overlays:
        - query: a.b.c.d
          value: {'foo': 'bar'}
          action: merge
          onMissing:
            action: "inject"
        - query: metadata.labels
          value: {'some': 'one'}
          action: merge
          onMissing:
            action: "inject"
        # demos multiple inject paths on missing
        - query: x.*
          value: {'x': 'x'}
          action: merge
          onMissing:
            action: "inject"
            injectPath:
              - x
              - y
              - z
  # demoing application of 'commonOverlays' without a 'overlays' or 'documents' key
  - name: "another file"
    path: "examples/manifests/another.yaml"
    # uncomment the following 3 lines to see this affect 2 of 3 docs in 'another.yaml' with commonOverlays {'cool_label': 'bro'}
    documents:
      - path: 0
      - path: 2
`

	// temporary delet this example either before PR or after PR
	s := `
--- 
yaml_files:
- name: "some arbitrary descriptor"
  path: "examples/manifests/test.yaml"
  overlays:
  - name: "delete all annotations"
    query: metadata.annotations
    value: "test"
    action: "delete"
`

	var test StaticYamlFiles
	// var errFile error
	// var data []byte
	// filename := fn

	// Catch No file passed
	// if data, errFile = ioutil.ReadFile(filename); errFile != nil {
	// 	return errFile
	// }

	err := yaml.Unmarshal([]byte(s), &test)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- t:\n%v\n\n", test)

	d, err := yaml.Marshal(&test)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

	var t models.Instructions

	err2 := yaml.Unmarshal([]byte(data), &t)
	if err2 != nil {
		log.Printf("error: %v", err2)
	}
	// fmt.Printf("--- t:\n%v\n\n", t)
	fmt.Printf("%+v", t.CommonOverlays[0].DocumentQuery[0])

	// fmt.Printf("File contents: %s", data)
	return nil
}
