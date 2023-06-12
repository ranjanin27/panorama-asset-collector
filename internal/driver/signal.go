// Copyright 2022 Hewlett Packard Enterprise Development LP

package driver

import (
	"os"
	"os/signal"
	"syscall"
)

// onlyOnceGuard is a channel that can be closed once and causes a panic
// if closed a second time.
var onlyOnceGuard = make(chan struct{})
var posixShutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

// RegisterShutdownSignalHandler calls the provided function when a shutdown
// signal (SIGINT or SIGTERM) is received.  A second shutdown signal causes
// the program to exit with status 1.  This function panics if called more
// than once.
func RegisterShutdownSignalHandler(handler func()) {
	close(onlyOnceGuard)

	//nolint:gomnd // Need room to receive two signals from a non-blocking sender.
	c := make(chan os.Signal, 2)
	signal.Notify(c, posixShutdownSignals...)
	go func() {
		<-c
		handler()
		<-c
		os.Exit(1)
	}()
}
