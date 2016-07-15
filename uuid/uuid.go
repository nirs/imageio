// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package uuid

import (
	"crypto/rand"
	"fmt"
)

type UUID [16]byte

// Uuid4 returns new UUID version 4
func Uuid4() (u UUID, err error) {
	if _, err = rand.Read(u[:]); err != nil {
		return
	}
	// Make it valid (see RFC 4122, section 4.4)
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	return
}

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}
