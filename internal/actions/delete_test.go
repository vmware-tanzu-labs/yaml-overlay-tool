// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/actions"
	"gopkg.in/yaml.v3"
)

func testInit() *yaml.Node {
	var data = `
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
		log.Fatalf("Error Unmarshalling: %s", err)
	}

	return &t
}

func TestDelete(t *testing.T) {
	testYaml := testInit()

	type args struct {
		root *yaml.Node
		path string
	}

	tests := []struct {
		name          string
		args          args
		wantErr       bool
		expectedValue string
	}{
		{
			name: "Delete Scalar Node",
			args: args{
				root: testYaml,
				path: "kind",
			},
			wantErr: false,
			expectedValue: `apiVersion: v1
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
			name: "Delete Map Node",
			args: args{
				root: testYaml,
				path: "metadata.annotations",
			},
			wantErr: false,
			expectedValue: `apiVersion: v1
metadata:
  name: bind-udp
  namespace: tanzu-dns
  labels:
    app.kubernetes.io/name: external-dns
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
			name: "Delete Seq Node",
			args: args{
				root: testYaml,
				path: "spec.ports[0]",
			},
			wantErr: false,
			expectedValue: `apiVersion: v1
metadata:
  name: bind-udp
  namespace: tanzu-dns
  labels:
    app.kubernetes.io/name: external-dns
spec:
  selector:
    app.kubernetes.io/name: external-dns
  type: LoadBalancer
  ports:
    - name: dns-tcp
      port: 53
      protocol: TCP
      targetPort: dns-tcp
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yp, _ := yamlpath.NewPath(tt.args.path)
			child, _ := yp.Find(tt.args.root)
			err := actions.Delete(tt.args.root, child[0], tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			buf := new(bytes.Buffer)
			ye := yaml.NewEncoder(buf)
			ye.SetIndent(2)
			ye.Encode(tt.args.root)
			if buf.String() != tt.expectedValue {
				t.Errorf("Delete() =\n%s, want \n%s", buf.String(), tt.expectedValue)
			}
		})
	}
}
