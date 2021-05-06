// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains tests for the yaml-overlay-tool 'merge' action

package actions_test

import (
	"bytes"
	"testing"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

func TestMerge(t *testing.T) {
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
			name: "Merge Scalar Node - string",
			args: args{
				query: "metadata.labels['app.kubernetes.io/name']",
				value: "es",
			},
			expectedValue: `apiVersion: v1
kind: Service
metadata:
  name: bind-udp
  namespace: tanzu-dns
  labels:
    app.kubernetes.io/name: external-dnses
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
			name: "Merge Scalar Node - integer",
			args: args{
				query: "spec.ports[1].port",
				value: "53",
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
      port: 106
      protocol: TCP
      targetPort: dns-tcp
`,
		},
		{
			name: "Merge Sequence Node (array)",
			args: args{
				query: "spec.ports",
				value: "- name: sequenceTest\n  port: 22\n  protocol: TCP\n  targetPort: sshd",
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
      port: 53
      protocol: TCP
      targetPort: dns-tcp
    - name: sequenceTest
      port: 22
      protocol: TCP
      targetPort: sshd
`,
		},
		{
			name: "Merge map",
			args: args{
				query: "metadata.labels",
				value: "foo: bar\nbar: foo\npotato: badayda",
			},
			expectedValue: `apiVersion: v1
kind: Service
metadata:
  name: bind-udp
  namespace: tanzu-dns
  labels:
    app.kubernetes.io/name: external-dns
    foo: bar
    bar: foo
    potato: badayda
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
			name: "Merge Scalar Node - string (with comments)",
			args: args{
				query: "metadata.annotations['service.beta.kubernetes.io/aws-load-balancer-type']",
				value: "# head\ns #line\n# foot",
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
	# head
    service.beta.kubernetes.io/aws-load-balancer-type: nlbs # line
	# foot
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
	}

	for _, tt := range tests {
		testYaml, val := testInit(tt.args.value)
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			yp, _ := yamlpath.NewPath(testCase.args.query)
			results, _ := yp.Find(testYaml)

			if err := actions.Merge(results[0], val.Content[0]); err != nil {
				t.Errorf("Encountered Error on merge action: %s", err)
			}

			buf := new(bytes.Buffer)
			ye := yaml.NewEncoder(buf)

			ye.SetIndent(2)

			if err := ye.Encode(testYaml); err != nil {
				t.Errorf("Encountered Error creating encoder: %s", err)
			}

			if buf.String() != testCase.expectedValue {
				t.Errorf("Merge() =\n%swant:\n%s", buf.String(), tt.expectedValue)
			}
		})
	}
}
