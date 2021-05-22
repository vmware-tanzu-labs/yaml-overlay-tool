// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions_test

import (
	"testing"

	"github.com/vmware-tanzu-labs/yaml-overlay-tool/internal/instructions"
)

func TestExecute(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg *instructions.Config
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test simple execution",
			args: args{
				cfg: &instructions.Config{
					InstructionsFile: "testdata/instructions.yaml",
					OutputDir:        "testdata/output/",
					Indent:           2,
				},
			},
			wantErr: false,
		},
		{
			name: "test diretory execution",
			args: args{
				cfg: &instructions.Config{
					InstructionsFile: "testdata/test.yaml",
					OutputDir:        "testdata/output/",
					Indent:           2,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := instructions.Execute(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Execute(%v) error = %v, wantErr %v", tt.args.cfg, err, tt.wantErr)
			}
		})
	}
}
