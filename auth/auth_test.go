// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package auth

import (
	"testing"
)

func TestMayReadNoAuth(t *testing.T) {
	u, _ := MayRead("3facfbc1", 1024)
	if u != nil {
		t.Fatal("Read allowed without a tikcet: %v", u)
	}
}

func TestMayWriteNoAuth(t *testing.T) {
	u, _ := MayWrite("3facfbc1", 1024)
	if u != nil {
		t.Fatal("Write allowed without a tikcet: %v", u)
	}
}

func TestAddRemove(t *testing.T) {
	ticket := &Ticket{
		Mode:    "r",
		Size:    1024,
		Timeout: 1,
		Url:     "file:///path",
		Uuid:    "3facfbc1",
	}
	err := Add(ticket)
	if err != nil {
		t.Fatal(err)
	}
	defer Remove(ticket.Uuid)

	u, err := MayRead(ticket.Uuid, 1024)
	if u == nil {
		t.Fatal("Auth not added: %v", err)
	}

	Remove(ticket.Uuid)
	u, err = MayRead(ticket.Uuid, 1024)
	if u != nil {
		t.Fatal("Auth not removed: %v", u)
	}
}

var mayRead = []*Ticket{
	{Mode: "rw", Size: 1024, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "r", Size: 1024, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "r", Size: 1025, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
}

var mayNotRead = []*Ticket{
	{Mode: "w", Size: 1024, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "r", Size: 1023, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "r", Size: 1024, Timeout: 0, Url: "file:///path", Uuid: "3facfbc1"},
}

var mayWrite = []*Ticket{
	{Mode: "rw", Size: 1024, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "w", Size: 1024, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "w", Size: 1025, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
}

var mayNotWrite = []*Ticket{
	{Mode: "r", Size: 1024, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "w", Size: 1023, Timeout: 1, Url: "file:///path", Uuid: "3facfbc1"},
	{Mode: "w", Size: 1024, Timeout: 0, Url: "file:///path", Uuid: "3facfbc1"},
}

func TestMayRead(t *testing.T) {
	for _, ticket := range mayRead {
		err := Add(ticket)
		if err != nil {
			t.Fatal(err)
		}
		defer Remove(ticket.Uuid)

		u, err := MayRead(ticket.Uuid, 1024)
		if err != nil {
			t.Errorf("Should allow read for %+v: %v", ticket, err)
			continue
		}
		if u.String() != ticket.Url {
			t.Errorf("Unexpected url %v", u)
		}
	}
}

func TestMayNotRead(t *testing.T) {
	for _, ticket := range mayNotRead {
		err := Add(ticket)
		if err != nil {
			t.Fatal(err)
		}
		defer Remove(ticket.Uuid)

		_, err = MayRead(ticket.Uuid, 1024)
		if err == nil {
			t.Errorf("Should not allow read for %+v", ticket)
		}
	}
}

func TestMayWrite(t *testing.T) {
	for _, ticket := range mayWrite {
		err := Add(ticket)
		if err != nil {
			t.Fatal(err)
		}
		defer Remove(ticket.Uuid)

		u, err := MayWrite(ticket.Uuid, 1024)
		if err != nil {
			t.Errorf("Should allow write for %+v: %v", ticket, err)
			continue
		}
		if u.String() != ticket.Url {
			t.Errorf("Unexpected url %v", u)
		}
	}
}

func TestMayNotWrite(t *testing.T) {
	for _, ticket := range mayNotWrite {
		err := Add(ticket)
		if err != nil {
			t.Fatal(err)
		}
		defer Remove(ticket.Uuid)

		_, err = MayWrite(ticket.Uuid, 1024)
		if err == nil {
			t.Errorf("Should not allow write for %+v", ticket)
		}
	}
}
