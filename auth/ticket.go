// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package auth

import (
	"encoding/json"
	"fmt"
)

type Bytes int64
type Seconds uint

type Ticket struct {
	Mode    string
	Size    Bytes
	Url     string
	Uuid    string
	Timeout Seconds
}

func ParseTicket(buf []byte) (t *Ticket, err error) {
	err = json.Unmarshal(buf, &t)
	if err != nil {
		return nil, fmt.Errorf("Invalid json: %v", err)
	}
	if !(t.Mode == "r" || t.Mode == "w" || t.Mode == "rw") {
		return nil, fmt.Errorf("Invalid mode: %v", t.Mode)
	}
	if t.Size <= 0 {
		return nil, fmt.Errorf("Size must be positive: %v", t.Size)
	}
	if t.Url == "" {
		return nil, fmt.Errorf("Url is required")
	}
	if t.Uuid == "" {
		return nil, fmt.Errorf("Uuid is required")
	}
	if t.Timeout == 0 {
		return nil, fmt.Errorf("Timeout is required")
	}
	return
}
