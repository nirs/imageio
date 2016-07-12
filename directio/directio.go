// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package directio

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// OpenFile opens a file with direct I/O enabled.
//
// Write and read to the file must use AlignedBuffer.
func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag|syscall.O_DIRECT, perm)
}

// Enable enables direct I/O on a file.
//
// When directio I/O is anabled, read and write must use AlignedBuffer.
func Enable(f *os.File) (err error) {
	flags, err := fcntl(f.Fd(), syscall.F_GETFL, 0)
	if err != nil {
		return
	}
	_, err = fcntl(f.Fd(), syscall.F_SETFL, flags|syscall.O_DIRECT)
	return
}

// Disable disables direct I/O on a file.
//
// When direct I/O is disabled, read and write can use unaligned buffer.
func Disable(f *os.File) (err error) {
	flags, err := fcntl(f.Fd(), syscall.F_GETFL, 0)
	if err != nil {
		return
	}
	_, err = fcntl(f.Fd(), syscall.F_SETFL, flags & ^syscall.O_DIRECT)
	return
}

func fcntl(fd uintptr, cmd uintptr, arg int) (val int, err error) {
	r, _, e := syscall.Syscall(syscall.SYS_FCNTL, fd, cmd, uintptr(arg))
	val = int(r)
	if e != 0 {
		err = e
	}
	return
}

// AlignedBuffer allocates aligned buffer.
//
// Allocate a buffer of at least size bytes, aligned to align bytes.
//
// size and align must be multiple of 512 bytes.
//
// Examples:
//  AlignedBuffer(1048576, 4096)
//
func AlignedBuffer(size int, align int) ([]byte, error) {
	if size%512 != 0 {
		return nil, fmt.Errorf("size must be mulitiple of 512: %v", size)
	}
	if align%512 != 0 {
		return nil, fmt.Errorf("align must be mulitiple of 512: %v", align)
	}
	buf := make([]byte, size+align)
	offset := 0
	remainder := int(uintptr(unsafe.Pointer(&buf[0])) & uintptr(align-1))
	if remainder != 0 {
		offset = align - remainder
	}
	return buf[offset : offset+size], nil
}
