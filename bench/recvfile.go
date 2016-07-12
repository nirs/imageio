// ovirt-imageio
// Copyright (C) 2016 Red Hat, Inc.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.

package main

import (
	"flag"
	"fmt"
	"os"
	"ovirt/imageio/fileio"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	MB = 1024 * 1024
)

var (
	offset   = flag.Int64("offset", 0, "offset in image (in MiB)")
	progress = flag.Bool("progress", false, "show copy progress (default false)")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: recvfile [options] PATH SIZE\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	path := args[0]

	size, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil || size < 0 {
		fmt.Printf("Error: invalid size: %s\n", size)
		os.Exit(2)
	}

	if *progress {
		pi := NewProgressIndicator(size*MB, time.Duration(100)*time.Millisecond)
		pi.Start()
		_, err = fileio.Receive(path, os.Stdin, size*MB, *offset*MB, pi)
		pi.Stop()
	} else {
		// Duplicated to avoid crash when sending (*T)(nil).
		_, err = fileio.Receive(path, os.Stdin, size*MB, *offset*MB, nil)
	}
	if err != nil {
		panic(err)
	}
}

type ProgressIndicator struct {
	value   int64 // Must be aligned to 8 byte to use atomic
	total   int64
	refresh time.Duration
	start   time.Time
	stopped chan bool
}

func NewProgressIndicator(total int64, refresh time.Duration) *ProgressIndicator {
	return &ProgressIndicator{
		total:   total,
		refresh: refresh,
		stopped: make(chan bool),
	}
}

func (p *ProgressIndicator) Start() {
	p.start = time.Now()
	go p.run()
}

func (p *ProgressIndicator) Stop() {
	p.stopped <- true
	p.update()
}

// Implement fileio.Progress interface
func (p *ProgressIndicator) Set(value int64) {
	atomic.StoreInt64(&p.value, value)
}

func (p *ProgressIndicator) run() {
	ticker := time.NewTicker(p.refresh)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			p.update()
		case <-p.stopped:
			return
		}
	}
}

func (p *ProgressIndicator) update() {
	elapsed := float64(time.Since(p.start)) / float64(time.Second)
	value := atomic.LoadInt64(&p.value)
	precent := fmt.Sprintf("%6.2f%%", float64(value)/float64(p.total)*100)
	rate := float64(value) / elapsed / float64(MB)
	fmt.Fprintf(os.Stdout, "[ %s ] %5.02f MiB/s %5.02fs\r", precent, rate, elapsed)
}
