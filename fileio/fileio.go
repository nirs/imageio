// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package fileio

import (
	"fmt"
	"io"
	"os"
	"ovirt/imageio/directio"
)

const (
	bufsize   = 8 * 1024 * 1024
	alignment = 4096
)

// Progress is an interface for reporting operation progress.
type Progress interface {
	Set(value int64)
}

// Receive copies size bytes from reader to path, staring at offset.
//
// Todo:
// - Check the option of Operation struct keeping the current progress, like
//   the Python version. Can be useful when we need to add cancling support.
// - There are too many parameters, would be nice to accept options struct
//   instead.
func Receive(path string, reader io.Reader, size int64, offset int64, progress Progress) (received int64, err error) {
	if size%512 != 0 {
		return 0, fmt.Errorf("size is not a multiple of 512 bytes: %v", size)
	}

	if offset%512 != 0 {
		return 0, fmt.Errorf("offset is not a multiple of 512 bytes: %v", offset)
	}

	file, err := directio.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer file.Close()

	if offset > 0 {
		if _, err = file.Seek(offset, os.SEEK_SET); err != nil {
			return
		}
	}

	buf, err := directio.AlignedBuffer(bufsize, alignment)
	if err != nil {
		return
	}

	// This is mostly like io.CopyBuffer, but:
	// - always read full blcoks, required for direct I/O
	// - read exactly size bytes instead of up to EOF
	//   must for http request, we don't want to read trailing bytes.
	//   This can be also implemented using http.MaxBytesReader, but
	//   it seems too complext for our needs.
	// - report progress after each write

	for received < size {
		b := buf
		todo := int(size - received)
		if todo < len(buf) {
			b = buf[:todo]
		}

		_, er := io.ReadFull(reader, b)
		if er != nil {
			err = er
			break
		}

		n, ew := file.Write(b)
		if n > 0 {
			received += int64(n)
			if progress != nil {
				progress.Set(received)
			}
		}
		if ew != nil {
			// file.Write handles EINTR and short writes; error means we
			// cannot write any more.
			err = ew
			break
		}
	}

	if se := file.Sync(); se != nil && err == nil {
		err = se
	}

	return
}
