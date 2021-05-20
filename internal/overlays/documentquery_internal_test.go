// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package overlays

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func testInit() (test *yaml.Node) {
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
  # add some fake boolean values for testing
  boolTest:
    case0: false
    case1: true
`

	var t yaml.Node

	err := yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("Error Unmarshalling test data: %s", err)
	}

	return &t
}

func TestDocumentQueries_checkQueries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		dq      DocumentQueries
		want    bool
		wantErr bool
	}{
		{
			name: "query should return true if it exists",
			dq: DocumentQueries{
				{
					Conditions: []*Condition{
						{
							Key: "kind",
							Value: yaml.Node{
								Kind:  yaml.ScalarNode,
								Tag:   "!!str",
								Value: `Service`,
							},
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name:    "query should return true if there are no conditions",
			dq:      DocumentQueries{},
			want:    false,
			wantErr: false,
		},
		{
			name: "query should return false it it does not exist",
			dq: DocumentQueries{
				{
					Conditions: []*Condition{
						{
							Key: "kind",
							Value: yaml.Node{
								Kind:  yaml.ScalarNode,
								Tag:   "!!str",
								Value: `Potato`,
							},
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			y := testInit()
			got, err := tt.dq.checkQueries(y)
			if (err != nil) != tt.wantErr {
				t.Errorf("DocumentQueries.checkQueries(%v) error = %v, wantErr %v", y, err, tt.wantErr)

				return
			}
			if got != tt.want {
				t.Errorf("DocumentQueries.checkQueries(%v) = %v, want %v", y, got, tt.want)
			}
		})
	}
}
