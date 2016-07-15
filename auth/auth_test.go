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

func TestParseTicket(t *testing.T) {
	text := `{
		"mode": "rw",
		"size": 1024,
		"timeout": 300,
		"url": "file:///path",
		"uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"
	}`
	ticket, err := ParseTicket([]byte(text))
	if err != nil {
		t.Fatal(err)
	}
	if ticket.Mode != "rw" {
		t.Fatalf("Unexpected mode: %+v", ticket)
	}
	if ticket.Size != 1024 {
		t.Fatalf("Unexpected size: %+v", ticket)
	}
	if ticket.Timeout != 300 {
		t.Fatalf("Unexpected timeout: %+v", ticket)
	}
	if ticket.Url != "file:///path" {
		t.Fatalf("Unexpected url: %+v", ticket)
	}
	if ticket.Uuid != "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2" {
		t.Fatalf("Unexpected uuid: %+v", ticket)
	}
}

var invalidTickets = []struct {
	desc string
	json string
}{
	{
		"Invalid mode",
		`{"mode": "x", "size": 1024, "timeout": 300, "url": "file:///path",
		  "uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"}`,
	},
	{
		"Missing mode",
		`{"size": 1024, "timeout": 300, "url": "file:///path",
		  "uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"}`,
	},
	{
		"Negative size",
		`{"mode": "rw", "size": -1, "timeout": 300, "url": "file:///path",
		  "uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"}`,
	},
	{
		"Missing size",
		`{"mode": "rw", "timeout": 300, "url": "file:///path",
		  "uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"}`,
	},
	{
		"Negative timeout",
		`{"mode": "rw", "size": 1024, "timeout": -1, "url": "file:///path",
		  "uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"}`,
	},
	{
		"Missing timeout",
		`{"mode": "rw", "size": 1024, "url": "file:///path",
		  "uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"}`,
	},
	{
		"Missing url",
		`{"mode": "rw", "size": 1024, "timeout": 300,
		  "uuid": "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"}`,
	},
	{
		"Missing uuid",
		`{"mode": "rw", "size": 1024, "timeout": 300, "url": "file:///path"}`,
	},
}

func TestInvalidTicket(t *testing.T) {
	for _, test := range invalidTickets {
		ticket, err := ParseTicket([]byte(test.json))
		if err == nil {
			t.Errorf("%s did not fail: %+v", test.desc, ticket)
		}
	}
}

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
