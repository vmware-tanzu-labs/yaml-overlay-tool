// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"reflect"
	"testing"
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
			want: map[string]interface{}{
				"too":  "bar",
				"yoo":  "car",
				"foo":  "you",
				"test": "one",
				"two":  "three",
			},
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
			want: map[string]interface{}{
				"too": "bar",
				"foo": "foo",
				"values": []interface{}{
					1,
					2,
					3,
				},
				"yoo": "bar",
			},
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
			want: map[string]interface{}{
				"too":  "bar",
				"yoo":  "czar",
				"test": "one",
				"two":  "three",
				"values": []interface{}{
					1,
					2,
					3,
				},
				"barf": "oo",
			},
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getValues(%v) = %+v, want %+v", tt.args.fileNames, got, tt.want)
			}
		})
	}
}
