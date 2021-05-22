// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains tests for the yaml-overlay-tool 'merge' action

package actions_test

import (
	"testing"
)

func TestMergeNode(t *testing.T) {
	t.Parallel()

	tests := testCases{
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`,
		},
		{
			name: "Merge map (inject line comment)",
			args: args{
				query: "metadata.labels",
				value: "foo: bar\nbar: foo\npotato: badayda # sometimes encountered in the northeast",
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
    potato: badayda # sometimes encountered in the northeast
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
			name: "Merge map (with line comment)",
			args: args{
				query: "metadata.labels",
				value: `app.kubernetes.io/name: "%s" # test`,
			},
			expectedValue: `apiVersion: v1
kind: Service
metadata:
  name: bind-udp
  namespace: tanzu-dns
  labels:
    app.kubernetes.io/name: external-dns # test
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
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
    - name: dns-udp-port # line
      # foot

      # head
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
	tests.runTests(t, "merge")
}
