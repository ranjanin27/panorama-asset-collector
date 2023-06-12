/*
 * (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 */

package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/handlers/kubernetes"
)

func mapRESTUrls(ctx context.Context) {
	// Panorama fleet hauler service metrics path for REST prometheus scraping
	router.GET("/metrics", func(c *gin.Context) {
		handlermetrics := promhttp.Handler()
		handlermetrics.ServeHTTP(c.Writer, c.Request)
	})

	// Healthz API for liveness and readiness
	healthHandler := kubernetes.NewKubeHandler(ctx, &kubernetes.HealthHandlerStruct{})
	router.GET("/healthz/liveness", healthHandler.GetHealth)
	router.GET("/healthz/readiness", healthHandler.GetReady)
}
