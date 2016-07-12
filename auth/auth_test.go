// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package auth_test

import (
	"ovirt/imageio/auth"
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
	ticket, err := auth.ParseTicket([]byte(text))
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
		ticket, err := auth.ParseTicket([]byte(test.json))
		if err == nil {
			t.Errorf("%s did not fail: %+v", test.desc, ticket)
		}
	}
}

var mayRead = []*auth.Ticket{
	{Mode: "rw", Size: 1024, Timeout: 1, Url: "file:///path"},
	{Mode: "r", Size: 1024, Timeout: 1, Url: "file:///path"},
	{Mode: "r", Size: 1025, Timeout: 1, Url: "file:///path"},
}

var mayNotRead = []*auth.Ticket{
	{Mode: "w", Size: 1024, Timeout: 1, Url: "file:///path"},
	{Mode: "r", Size: 1023, Timeout: 1, Url: "file:///path"},
	{Mode: "r", Size: 1024, Timeout: 0, Url: "file:///path"},
}

var mayWrite = []*auth.Ticket{
	{Mode: "rw", Size: 1024, Timeout: 1, Url: "file:///path"},
	{Mode: "w", Size: 1024, Timeout: 1, Url: "file:///path"},
	{Mode: "w", Size: 1025, Timeout: 1, Url: "file:///path"},
}

var mayNotWrite = []*auth.Ticket{
	{Mode: "r", Size: 1024, Timeout: 1, Url: "file:///path"},
	{Mode: "w", Size: 1023, Timeout: 1, Url: "file:///path"},
	{Mode: "w", Size: 1024, Timeout: 0, Url: "file:///path"},
}

func TestMayRead(t *testing.T) {
	for _, ticket := range mayRead {
		may, err := auth.NewAuth(ticket)
		if err != nil {
			t.Errorf("Cannot create auth for %+v: %v", ticket, err)
			continue
		}
		u, err := may.Read(1024)
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
		may, err := auth.NewAuth(ticket)
		if err != nil {
			t.Errorf("Cannot create auth for %+v: %v", ticket, err)
			continue
		}
		_, err = may.Read(1024)
		if err == nil {
			t.Errorf("Should not allow read for %+v", ticket)
		}
	}
}

func TestMayWrite(t *testing.T) {
	for _, ticket := range mayWrite {
		may, err := auth.NewAuth(ticket)
		if err != nil {
			t.Errorf("Cannot create auth for %+v: %v", ticket, err)
			continue
		}
		u, err := may.Write(1024)
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
		may, err := auth.NewAuth(ticket)
		if err != nil {
			t.Errorf("Cannot create auth for %+v: %v", ticket, err)
			continue
		}
		_, err = may.Write(1024)
		if err == nil {
			t.Errorf("Should not allow write for %+v", ticket)
		}
	}
}

func TestGetNoAuth(t *testing.T) {
	may, _ := auth.Get("ticket-uuid")
	if may != nil {
		t.Fatal("Unexpected auth: %v", may)
	}
}

func TestAdd(t *testing.T) {
	u := "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"
	ticket := &auth.Ticket{
		Mode:    "r",
		Size:    1024,
		Timeout: 1,
		Url:     "file:///path",
		Uuid:    u,
	}
	err := auth.Add(ticket)
	if err != nil {
		t.Fatal(err)
	}
	may, err := auth.Get(u)
	if may == nil || err != nil {
		t.Fatalf("Expected auth, got nil: %v", err)
	}
}

func TestRemove(t *testing.T) {
	u := "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"
	ticket := &auth.Ticket{
		Mode:    "r",
		Size:    1024,
		Timeout: 1,
		Url:     "file:///path",
		Uuid:    u,
	}
	err := auth.Add(ticket)
	if err != nil {
		t.Fatal(err)
	}
	auth.Remove(u)
	may, _ := auth.Get(u)
	if may != nil {
		t.Fatalf("Unexpected auth: %v", may)
	}
}
