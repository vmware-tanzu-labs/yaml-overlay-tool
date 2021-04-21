// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"bufio"
	"io"
	"os"
)

func ReadStream(fileName string) (io.Reader, error) {
	if fileName == "-" {
		return bufio.NewReader(os.Stdin), nil
	}

	return os.Open(fileName)
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Error("error closing file!: %s", err)
	}
}
