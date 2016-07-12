// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package images_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"ovirt/imageio/auth"
	"ovirt/imageio/images"
	"ovirt/imageio/testutil"
	"testing"
)

func TestGetNotFound(t *testing.T) {
	err := images.Start("localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer images.Stop()

	resp, err := request("GET", "/images/no-such-ticket", nil)
	if resp == nil {
		t.Fatalf("Request failed: err=%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected %v, got %v", http.StatusNotFound, resp.StatusCode)
	}
}

func TestPutNotFound(t *testing.T) {
	err := images.Start("localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer images.Stop()

	resp, err := request("PUT", "/images/no-such-ticket", nil)
	if resp == nil {
		t.Fatalf("Request failed: err=%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Expected %v, got %v", http.StatusNotFound, resp.StatusCode)
	}
}

func TestPut(t *testing.T) {
	err := images.Start("localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer images.Stop()

	const size = 1024

	path, err := testutil.CreateFile(size)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	u := "3facfbc1-68e0-4b77-b0c6-87e66fcabcc2"
	ticket := &auth.Ticket{
		Mode:    "rw",
		Size:    size,
		Timeout: 10,
		Url:     "file://" + path,
		Uuid:    u,
	}
	err = auth.Add(ticket)
	if err != nil {
		t.Fatal(err)
	}
	defer auth.Remove(u)

	buf := testutil.Buffer(size)
	resp, err := request("PUT", "/images/"+u, buf)
	if resp == nil {
		t.Fatalf("Request failed: err=%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected %v, got %v", http.StatusOK, resp.StatusCode)
	}

	content, err := ioutil.ReadFile(path)
	if !bytes.Equal(content, buf) {
		t.Fatalf("Expected %v, got %v", buf, content)
	}
}

func TestAlreadyRunning(t *testing.T) {
	err := images.Stop()
	if err == nil {
		t.Fatal("Stop did not fail on stopped server")
	}
}

func TestNotRunning(t *testing.T) {
	err := images.Start("localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer images.Stop()

	err = images.Start("ignored")
	if err == nil {
		t.Fatal("Start did not fail on running server")
	}
}

// request sends http request ot the images server
func request(method string, path string, buf []byte) (resp *http.Response, err error) {
	url := fmt.Sprintf("http://%s%s", images.Addr(), path)
	body := bytes.NewReader(buf)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	if body != nil {
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(buf)))
		req.Header.Set("Content-Type", "application/octet-stream")
	}
	return http.DefaultClient.Do(req)
}
