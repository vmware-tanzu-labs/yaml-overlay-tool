// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains common test functions and types for yaml-overlay-tool actions

package actions_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

type args struct {
	query string
	value string
}

type testCase struct {
	name          string
	args          args
	expectedValue string
}

type testCases []testCase

func testInit(v string) (orig, test *yaml.Node) {
	data := `
apiVersion: v1
kind: Service
metadata:
  name: bind-udp
  namespace: tanzu-dns
  labels:
    app.kubernetes.io/name: external-dns
  annotations:
    # NOTE: this only works on 1.19.1+vmware.1+, but not prior
    ## This annotation will be ignored on other cloud providers
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
spec:
  selector:
    app.kubernetes.io/name: external-dns
  type: LoadBalancer
  ports:
    - name: dns-udp
      port: 53
      protocol: UDP
      targetPort: dns-udp
    - name: dns-tcp
      port: 53
      protocol: TCP
      targetPort: dns-tcp
`

	var t yaml.Node

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("Error Unmarshalling test data: %s", err)
	}

	var val yaml.Node

	if v != "" {
		e := yaml.Unmarshal([]byte(v), &val)
		if e != nil {
			log.Fatalf("Error Unmarshalling test value: %s", err)
		}
	}

	return &t, &val
}

func (tst testCases) runTests(t *testing.T, a string) {
	t.Helper()

	for _, tt := range tst {
		testYaml, val := testInit(tt.args.value)
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			yp, _ := yamlpath.NewPath(testCase.args.query)
			results, _ := yp.Find(testYaml)

			testAction(t, a, results[0], testYaml, val.Content...)

			buf := new(bytes.Buffer)
			ye := yaml.NewEncoder(buf)

			ye.SetIndent(2)

			if err := ye.Encode(testYaml); err != nil {
				t.Errorf("Encountered Error creating encoder: %s", err)
			}

			if buf.String() != testCase.expectedValue {
				t.Errorf("%s() =\n%swant:\n%s", a, buf.String(), tt.expectedValue)
			}
		})
	}
}

func testAction(t *testing.T, action string, result, testYaml *yaml.Node, testValue ...*yaml.Node) {
	t.Helper()

	switch action {
	case "merge":
		if err := actions.Merge(result, testValue[0]); err != nil {
			t.Errorf("Encountered Error on merge action: %s", err)
		}
	case "replace":
		if err := actions.Replace(result, testValue[0]); err != nil {
			t.Errorf("Encountered Error on replace action: %s", err)
		}
	case "format":
		if err := actions.Format(result, testValue[0]); err != nil {
			t.Errorf("Encountered Error on format action: %s", err)
		}
	case "delete":
		actions.Delete(testYaml, result)
	}
}
