// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

func renderInstructionsTemplate(fileName string, values interface{}) (io.Reader, error) {
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to process template for instructions file %s, %w", fileName, err)
	}

	tmpl = tmpl.Option("missingkey=zero")

	var b bytes.Buffer

	if err := tmpl.Execute(&b, values); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &b, nil
}
