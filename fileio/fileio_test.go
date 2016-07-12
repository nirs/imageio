// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package fileio_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"ovirt/imageio/fileio"
	"ovirt/imageio/testutil"
	"testing"
)

func TestReceiveFull(t *testing.T) {
	const size = 1024 * 1234
	path, err := testutil.CreateFile(size)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	buf := testutil.Buffer(size)
	reader := bytes.NewReader(buf)
	n, err := fileio.Receive(path, reader, size, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != size {
		t.Fatalf("Received %v bytes, expected %v bytes", n, size)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, buf) {
		t.Fatalf("Expected %v, got %v", buf, content)
	}
}

func TestReceiveRandom(t *testing.T) {
	const size = 1024 * 1234
	path, err := testutil.CreateFile(size)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	buf := testutil.Buffer(size)
	reader := testutil.RandomReader(bytes.NewReader(buf))
	n, err := fileio.Receive(path, reader, size, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != size {
		t.Fatalf("Received %v bytes, expected %v bytes", n, size)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(content, buf) {
		t.Fatalf("Expected %v, got %v", buf, content)
	}
}

func TestReceiveOffset(t *testing.T) {
	const filesize = 512 * 3
	const bufsize = 512
	const offset = 512
	path, err := testutil.CreateFile(filesize)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	buf := testutil.Buffer(bufsize)
	reader := bytes.NewReader(buf)
	n, err := fileio.Receive(path, reader, bufsize, offset, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != bufsize {
		t.Fatalf("Received %v bytes, expected %v bytes", n, bufsize)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	empty := make([]byte, 512)
	block1 := content[:512]
	if !bytes.Equal(block1, empty) {
		t.Fatalf("Expected %v, got %v", empty, block1)
	}

	block2 := content[512:1024]
	if !bytes.Equal(block2, buf) {
		t.Fatalf("Expected %v, got %v", buf, block2)
	}

	block3 := content[1024:]
	if !bytes.Equal(block3, empty) {
		t.Fatalf("Expected %v, got %v", empty, block3)
	}
}

func TestReceiveUnalignedSize(t *testing.T) {
	const size = 511
	path, err := testutil.CreateFile(size)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	buf := testutil.Buffer(size)
	reader := bytes.NewReader(buf)
	n, err := fileio.Receive(path, reader, size, 0, nil)
	if n != 0 || err == nil {
		t.Fatalf("Call did not fail: n=%v, err=%v", n, err)
	}
}

func TestReceiveUnalignedOffset(t *testing.T) {
	path, err := testutil.CreateFile(1024)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	buf := testutil.Buffer(512)
	reader := bytes.NewReader(buf)
	n, err := fileio.Receive(path, reader, 512, 511, nil)
	if n != 0 || err == nil {
		t.Fatalf("Call did not fail: n=%v, err=%v", n, err)
	}
}
