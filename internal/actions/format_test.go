// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains tests for the yaml-overlay-tool 'format' action

package actions_test

import (
	"bytes"
	"testing"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

func TestFormat(t *testing.T) {
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
			name: "Format Scalar Node (only type accepted for format action)",
			args: args{
				query: "kind",
				value: "My%s",
			},
			expectedValue: `apiVersion: v1
kind: MyService
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
`,
		},
		{
			name: "Format Scalar Node Key",
			args: args{
				query: "spec.ports~",
				value: "s%s",
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
  sports:
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
			// Found that the '%s' cannot be unquoted at line start, or yaml cannot unmarshal
			// The quotation happens below in value via escaped double-quotes
			name: "Format Scalar Node in Array with head, line, foot comments",
			args: args{
				query: "spec.ports[0].name",
				value: "# head\n\"%s-port\" # line\n# foot",
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
	# head
    - name: dns-udp-port # line
	# foot
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

			if err := actions.Format(results[0], val.Content[0]); err != nil {
				t.Errorf("Encountered Error on format action: %s", err)
			}

			buf := new(bytes.Buffer)
			ye := yaml.NewEncoder(buf)

			ye.SetIndent(2)

			if err := ye.Encode(testYaml); err != nil {
				t.Errorf("Encountered Error creating encoder: %s", err)
			}

			if buf.String() != testCase.expectedValue {
				t.Errorf("Format() =\n%swant:\n%s", buf.String(), tt.expectedValue)
			}
		})
	}
}
