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

// Uuid4 returns new uuid version 4 string
func Uuid4() (u string, err error) {
	b := make([]byte, 16)
	if _, err = rand.Read(b); err != nil {
		return
	}
	// Make it valid (see RFC 4122, section 4.4)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	u = fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
