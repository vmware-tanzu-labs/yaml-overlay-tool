// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains tests for the yaml-overlay-tool 'replace' action

package actions_test

import (
	"testing"
)

func TestReplaceNode(t *testing.T) {
	t.Parallel()

	tests := testCases{
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`,
		},
	}

	tests.runTests(t, "replace")
}
