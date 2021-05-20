// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package actions

import (
	"fmt"
	"strings"
)

func CondSprintf(format string, v ...interface{}) string {
	v = append(v, "")
	format = strings.Replace(format, "%v", "%[1]v", -1)
	format = strings.Replace(format, "%l", "%[2]s", -1)
	format = strings.Replace(format, "%f", "%[3]s", -1)
	format = strings.Replace(format, "%h", "%[4]s", -1)

	format += fmt.Sprint("%[", len(v), "]s")

	return fmt.Sprintf(format, v...)
}
