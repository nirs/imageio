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
	"net/url"
	"strings"
	"sync"
	"time"
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

// Auth provide authorization based on ticket and creation time
type Auth struct {
	ticket  *Ticket
	expires time.Time
	url     *url.URL
}

var supportedSchemes = map[string]bool{"file": true}

// NewAuth creates new Auth from ticket, valid for t.Timeout seconds.
func NewAuth(t *Ticket) (*Auth, error) {
	u, err := url.Parse(t.Url)
	if err != nil {
		return nil, fmt.Errorf("Invalid url: %v: %v", t.Url, err)
	}
	if !supportedSchemes[u.Scheme] {
		return nil, fmt.Errorf("Unsupported scheme: %v", u.Scheme)
	}
	expires := time.Now().Add(time.Duration(t.Timeout) * time.Second)
	return &Auth{t, expires, u}, nil
}

// Read checks if caller may read up to size bytes, and return a url that the
// caller may read from, or an error describing why the operation is forbidden.
func (a *Auth) Read(size int64) (*url.URL, error) {
	return a.check("r", size)
}

// Write checks if caller may write up to size bytes, and return a url that the
// caller may write to, or an error describing why the operation is forbidden.
func (a *Auth) Write(size int64) (*url.URL, error) {
	return a.check("w", size)
}

func (a *Auth) check(mode string, size int64) (*url.URL, error) {
	if !strings.Contains(a.ticket.Mode, mode) {
		return nil, fmt.Errorf("Operation not allowed: %v", mode)
	}
	if size > int64(a.ticket.Size) {
		return nil, fmt.Errorf("Size out of range: %v", size)
	}
	if time.Now().After(a.expires) {
		return nil, fmt.Errorf("Ticket expired at %s", a.expires)
	}
	return a.url, nil
}

// Authorizations are accessed by multiple webserver goroutines /tickets/
// requests are adding and removing, and /images/ requests are getting.  We can
// use channels for synchronization, but single mutex seems simpler.  Extending
// auth may need additional locking so concurrent calls to Auth.check will see
// the correct expires time.

var (
	authorization = map[string]*Auth{}
	mutex         = sync.Mutex{}
)

// Returns Auth for ticket uuid, or error.
func Get(u string) (a *Auth, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	a = authorization[u]
	if a == nil {
		err = fmt.Errorf("No auth for %v", u)
	}
	return
}

// Add adds Auth for ticket
func Add(t *Ticket) (err error) {
	a, err := NewAuth(t)
	if err != nil {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	authorization[t.Uuid] = a
	return
}

// Remove removes Auth for u
func Remove(u string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(authorization, u)
	// TODO: cancel tasks authorized by u
}
