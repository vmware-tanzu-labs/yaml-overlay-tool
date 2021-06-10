// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package instructions

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func renderInstructionsTemplate(fileName string, values interface{}) (io.Reader, error) {
	tpl, err := template.ParseFiles(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to process template for instructions file %s, %w", fileName, err)
	}

	tpl = tpl.Funcs(template.FuncMap(sprig.FuncMap())).Option("missingkey=zero")

	var b bytes.Buffer

	if err := tpl.Execute(&b, values); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &b, nil
}
