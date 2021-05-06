// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains tests for the yaml-overlay-tool 'replace' action

package actions_test

import (
	"bytes"
	"testing"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

func TestReplace(t *testing.T) {
	t.Parallel()

	type args struct {
		query string
		value string
	}

	tests := []struct {
		name          string
		args          args
		expectedValue string
	}{
		{
			name: "Replace Scalar Node (string)",
			args: args{
				query: "metadata.labels['app.kubernetes.io/name']",
				value: "rpk-dns",
			},
			expectedValue: `apiVersion: v1
kind: Service
metadata:
  name: bind-udp
  namespace: tanzu-dns
  labels:
    app.kubernetes.io/name: rpk-dns
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
`,
		},
		{
			name: "Replace value in array (replace integer with string)",
			args: args{
				query: "spec.ports[1].port",
				value: "953",
			},
			expectedValue: `apiVersion: v1
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
      port: 953
      protocol: TCP
      targetPort: dns-tcp
`,
		},
		{
			name: "Replace entire array",
			args: args{
				query: "spec.ports",
				value: `- name: fake-data
  port: 9999
  protocol: TCP
  targetPort: my-fake-port
- name: fake-data2
  port: 1111
  protocol: UDP
  targetPort: another-fake-port
`,
			},
			expectedValue: `apiVersion: v1
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
    - name: fake-data
      port: 9999
      protocol: TCP
      targetPort: my-fake-port
    - name: fake-data2
      port: 1111
      protocol: UDP
      targetPort: another-fake-port
`,
		},
	}

	for _, tt := range tests {
		testYaml, val := testInit(tt.args.value)
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			yp, _ := yamlpath.NewPath(testCase.args.query)
			results, _ := yp.Find(testYaml)

			if err := actions.Replace(results[0], val.Content[0]); err != nil {
				t.Errorf("Encountered Error on replace action: %s", err)
			}

			buf := new(bytes.Buffer)
			ye := yaml.NewEncoder(buf)

			ye.SetIndent(2)

			if err := ye.Encode(testYaml); err != nil {
				t.Errorf("Encountered Error creating encoder: %s", err)
			}

			if buf.String() != testCase.expectedValue {
				t.Errorf("Replace() =\n%swant:\n%s", buf.String(), tt.expectedValue)
			}
		})
	}
}
