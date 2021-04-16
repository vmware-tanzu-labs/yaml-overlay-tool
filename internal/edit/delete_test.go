// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package edit_test

import (
	"log"
	"testing"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/edit"
	"gopkg.in/yaml.v3"
)

func testInit() *yaml.Node {
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

	var t yaml.Node

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatal("Error Unmarshalling")
	}

	return &t
}

func TestDeleteSeqNode(t *testing.T) {
	y := testInit()

	t.Run("Delete Test", func(t *testing.T) {
		edit.DeleteSeqNode(y.Content[0].Content[3], "name", "another file")
		o, _ := yaml.Marshal(y)
		ll := len(o)
		if ll != 3936 {
			t.Errorf("Document does not match expected length, got %d; want 3936", ll)
		}
	})
}
