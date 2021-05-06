// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

// contains common test functions for yaml-overlay-tool actions

package actions_test

import (
	"log"

	"gopkg.in/yaml.v3"
)

func testInit(v string) (*yaml.Node, *yaml.Node) {
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
