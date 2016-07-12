// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package testutil_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"ovirt/imageio/testutil"
	"strings"
	"testing"
)

func TestRandomReader(t *testing.T) {
	src := "Text to copy to writer"
	reader := testutil.RandomReader(strings.NewReader(src))
	writer := &bytes.Buffer{}
	_, err := io.Copy(writer, reader)
	if err != nil {
		t.Fatal(err)
	}
	dst := writer.String()
	if src != dst {
		t.Fatalf("Expected %#v, got %#v", src, dst)
	}
}

func TestRandomReaderOne(t *testing.T) {
	reader := testutil.RandomReader(strings.NewReader("123456789"))
	buf := make([]byte, 1)
	n, err := reader.Read(buf)
	if n != 1 {
		t.Fatalf("Read to one byte buffer failed, n=%v, err=%v", n, err)
	}
}

func TestRandomReaderEmpty(t *testing.T) {
	reader := testutil.RandomReader(strings.NewReader("123456789"))
	buf := make([]byte, 0)
	n, err := reader.Read(buf)
	if n != 0 {
		t.Fatalf("Read to empty buffer returned, n=%v, err=%v", n, err)
	}
}

func TestCreateFile(t *testing.T) {
	size := 1024 * 1024
	path, err := testutil.CreateFile(size)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(buf) != size {
		t.Fatalf("Read %v bytes, expected %v", len(buf), size)
	}
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			t.Fatal("Uexpected byte: %v", buf[i])
		}
	}
}

func TestBufferStart(t *testing.T) {
	buf := testutil.Buffer(4)
	expected := []byte{0, 1, 2, 3}
	if !bytes.Equal(buf, expected) {
		t.Fatalf("Expected %v, got %v", expected, buf)
	}
}

func TestBufferRepeat(t *testing.T) {
	buf := testutil.Buffer(255 + 4)[255:]
	expected := []byte{255, 0, 1, 2}
	if !bytes.Equal(buf, expected) {
		t.Fatalf("Expected %v, got %v", expected, buf)
	}
}
