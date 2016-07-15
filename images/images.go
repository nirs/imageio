// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package images

import (
	"fmt"
	"net"
	"net/http"
	"ovirt/imageio/auth"
	"ovirt/imageio/fileio"
)

const (
	ROOT = "/images/"
)

var (
	listener net.Listener
)

// Start starts the images web server.
func Start(addr string) (err error) {
	if listener != nil {
		return fmt.Errorf("Already started")
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	listener = ln

	mux := http.NewServeMux()
	mux.HandleFunc(ROOT, handle)
	server := &http.Server{Handler: mux}

	go server.Serve(listener)
	return
}

// Stop stops the images web server.
//
// This does not effect ongoing requests, and does not wait for their
// completion.
func Stop() error {
	if listener == nil {
		return fmt.Errorf("Not running")
	}
	// The documentaion is not clear about the semantics of close error.
	// Looking at the implementaiton, it seems that this error is very unlikely
	// as nobody else is using this listener.
	ln := listener
	listener = nil
	return ln.Close()
}

// Addr returns the address the server is listening on. For testing a server on
// a random port.
func Addr() string {
	return listener.Addr().String()
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PUT":
		put(w, r)
	case "GET":
		http.NotFound(w, r)
	default:
		http.Error(w, "You are not allowed to "+r.Method, http.StatusMethodNotAllowed)
		return
	}
}

func put(w http.ResponseWriter, r *http.Request) {
	ticketUuid := r.URL.Path[len(ROOT):]
	url, err := auth.MayWrite(ticketUuid, r.ContentLength)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_, err = fileio.Receive(url.Path, r.Body, r.ContentLength, 0, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
