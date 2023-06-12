// Copyright (c) 2022 Hewlett Packard Enterprise Development LP

package driver

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
)

func TestRegisterShutdownSignalHandler(t *testing.T) {
	RegisterShutdownSignalHandler(func() {})
	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt)
	signal.Notify(sigc, syscall.SIGTERM)
}
