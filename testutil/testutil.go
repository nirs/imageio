// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package testutil

import (
	"io"
	"io/ioutil"
	"math/rand"
)

// RandomReader returns random amount of bytes on each read.
//
// This reader is useful for testing coping to files using direct I/O,
// requiring aligned writes.
func RandomReader(r io.Reader) io.Reader {
	return &randomReader{r}
}

type randomReader struct {
	r io.Reader
}

func (r *randomReader) Read(p []byte) (int, error) {
	if len(p) > 1 {
		n := rand.Intn(len(p))
		p = p[0:n]
	}
	return r.r.Read(p)
}

// CreateFile creates a temporary file of size bytes full of zeros.
//
// Caller is responsible for removing the temporary file.
func CreateFile(size int) (string, error) {
	file, err := ioutil.TempFile("/var/tmp", "testutil.")
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.WriteAt([]byte{0}, int64(size-1))
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func Buffer(size int) []byte {
	buf := make([]byte, size)
	for i := 0; i < size; i++ {
		buf[i] = byte(i & 0xFF)
	}
	return buf
}
