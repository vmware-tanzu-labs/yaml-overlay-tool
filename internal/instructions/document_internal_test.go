// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"testing"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/overlays"
)

func TestDocument_checkDocumentIndex(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name     string
		Path     int
		Overlays []*overlays.Overlay
	}

	type args struct {
		docIndex int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "different indexes should return false",
			fields: fields{
				Name:     "test 1",
				Path:     0,
				Overlays: nil,
			},
			args: args{
				docIndex: 1,
			},
			want: false,
		},
		{
			name: "same indexes should return true",
			fields: fields{
				Name:     "test 1",
				Path:     0,
				Overlays: nil,
			},
			args: args{
				docIndex: 0,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := &Document{
				Name:     tt.fields.Name,
				Path:     tt.fields.Path,
				Overlays: tt.fields.Overlays,
			}
			if got := d.checkDocumentIndex(tt.args.docIndex); got != tt.want {
				t.Errorf("Document.checkDocumentIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
