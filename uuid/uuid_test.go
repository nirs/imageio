// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package uuid_test

import (
	"ovirt/imageio/uuid"
	"regexp"
	"testing"
)

func TestUuid4Format(t *testing.T) {
	u, err := uuid.Uuid4()
	if err != nil {
		t.Fatal(err)
	}
	// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_.28random.29
	pat := "[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}"
	if _, err := regexp.Match(pat, []byte(u)); err != nil {
		t.Fatalf("No match for %v: %v", u, err)
	}
}

func TestUuid4Uniq(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 10000; i++ {
		u, err := uuid.Uuid4()
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
		uuid.Uuid4()
	}
}
