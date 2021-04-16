// Copyright 2021 VMware, Inc.
// SPDX-License-Identifier: MIT

package lib

import (
	"bufio"
	"io"
	"os"

	"log"
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
		log.Printf("error closing file!: %s", err)
	}
}
