// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package uuid

import (
	"regexp"
	"testing"
)

func TestUuid4String(t *testing.T) {
	u, err := Uuid4()
	if err != nil {
		t.Fatal(err)
	}
	us := u.String()
	// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
	pat := "[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}"
	if matched, err := regexp.Match(pat, []byte(us)); !matched {
		t.Fatalf("No match for %v: %v", us, err)
	}
}

func TestUuid4Uniq(t *testing.T) {
	seen := make(map[UUID]bool)
	for i := 0; i < 10000; i++ {
		u, err := Uuid4()
		if err != nil {
			t.Fatal(err)
		}
		if seen[u] {
			t.Fatalf("Seen %v", u)
		}
		seen[u] = true
	}
}

func BenchmarkUuid4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Uuid4()
	}
}

func BenchmarkUuid4String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		u, _ := Uuid4()
		u.String()
	}
}
