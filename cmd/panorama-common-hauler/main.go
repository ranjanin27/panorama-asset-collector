// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package main

import (
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/app"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
)

// filled in by ldflags during build
var (
	Version string
	logger  = logging.GetLogger()
)

func main() {
	logger.Infof("Starting Panorama Common Hauler micro service with version %s", Version)
	app.StartApplication(Version)
}
