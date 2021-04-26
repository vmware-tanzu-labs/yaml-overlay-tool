// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadStream(fileName string) (io.Reader, error) {
	if fileName == "-" {
		return bufio.NewReader(os.Stdin), nil
	}

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %w", fileName, err)
	}

	return file, nil
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Error("error closing file!: %s", err)
	}
}
