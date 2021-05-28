// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func Test_getValues(t *testing.T) {
	t.Parallel()

	type args struct {
		fileNames []string
	}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test merging simple values",
			args: args{
				fileNames: []string{
					"testdata/values/defaults.yaml",
					"testdata/values/on_premise.yml",
				},
			},
			want: `foo: you
test: one
too: bar
two: three
yoo: car
`,
			wantErr: false,
		},
		{
			name: "test merge of several override values",
			args: args{
				fileNames: []string{
					"testdata/values/defaults.yaml",
					"testdata/values/site_b.yml",
				},
			},
			want: `foo: foo
too: bar
values:
    - 1
    - 2
    - 3
yoo: bar
`,
			wantErr: false,
		},
		{
			name: "multi-file complex merge",
			args: args{
				fileNames: []string{
					"testdata/values/defaults.yaml",
					"testdata/values/on_premise.yml",
					"testdata/values/site_b.yml",
					"testdata/values/site_a.yaml",
				},
			},
			want: `foo:
    bar:
        potato:
            badayda:
                - 1
                - cheese: burger
                - soup:
                    - minestrone
                    - chicken noodle
                    - clam chowder
test: one
too: bar
two: three
values:
    - 1
    - 2
    - 3
yoo: czar
zarf: oo
`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getValues(tt.args.fileNames)
			if (err != nil) != tt.wantErr {
				t.Errorf("getValues(%v) error = %v, wantErr %v", tt.args.fileNames, err, tt.wantErr)

				return
			}
			b, _ := yaml.Marshal(got)
			if !reflect.DeepEqual(string(b), tt.want) {
				t.Errorf("getValues(%v) = %v, want %v", tt.args.fileNames, string(b), tt.want)
			}
		})
	}
}
