// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains tests for the yaml-overlay-tool 'math' action

package actions_test

import (
	"testing"
)

func TestMathNode(t *testing.T) {
	t.Parallel()

	tests := testCases{
		{
			name: "Math Scalar Node - string",
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`,
		},
		{
			name: "Math Scalar Node - integer",
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`,
		},
		{
			name: "Math Scalar Node - true boolean with false boolean",
			args: args{
				query: "spec.boolTest.case0",
				value: "true",
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`,
		},
		{
			name: "Math Scalar Node - false boolean with false boolean",
			args: args{
				query: "spec.boolTest.case0",
				value: "false",
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`,
		},
		{
			name: "Math Scalar Node - true boolean with true boolean",
			args: args{
				query: "spec.boolTest.case1",
				value: "true",
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`,
		},
		{
			name: "Math Scalar Node - false boolean with true boolean",
			args: args{
				query: "spec.boolTest.case1",
				value: "false",
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: false
`,
		},
		{
			name: "Math Scalar Node - string (with comments)",
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
    service.beta.kubernetes.io/aws-load-balancer-type: nlbs #line
    # foot

    # head
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
	}

	tests.runTests(t, "math")
}
