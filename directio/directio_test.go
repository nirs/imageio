// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package directio

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestBadSize(t *testing.T) {
	_, err := AlignedBuffer(100, 4096)
	if err == nil {
		t.Fatal("Expected error with unaligned size")
	}
}

func TestBadAlign(t *testing.T) {
	_, err := AlignedBuffer(1024, 100)
	if err == nil {
		t.Fatal("Expected error with unaligned align")
	}
}

func TestOpenFile(t *testing.T) {
	dir, err := ioutil.TempDir("/var/tmp", "directio.")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	name := filepath.Join(dir, "file")
	file, err := OpenFile(name, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	buf, err := AlignedBuffer(1024, 4096)
	if err != nil {
		t.Fatal(err)
	}

	copy(buf, "testing 1 2 3...")

	n, err := file.Write(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(buf) {
		t.Fatalf("Wrote %v bytes, expected %v", n, len(buf))
	}

	file.Seek(0, 0)
	n, err = file.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(buf) {
		t.Fatalf("Read %v bytes, expected %v", n, len(buf))
	}
}
