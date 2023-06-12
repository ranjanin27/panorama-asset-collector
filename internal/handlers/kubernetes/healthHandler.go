// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package kubernetes

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Status = "status"
	Ready  = "Ready"
	Alive  = "Alive"
)

type KubeHandlerInterface interface {
	// GetHealth returns a status 200 for liveness
	GetHealth(c *gin.Context)
	// GetReady returns a status 200 for readiness
	GetReady(c *gin.Context)
}
type kubeHandlerStruct struct {
	ctx                 context.Context
	healthHandlerStruct HealthHandlerInterface
}

type HealthHandlerInterface interface {
	GetKafkaHealth() error
}
type HealthHandlerStruct struct{}

func NewKubeHandler(ctx context.Context, healthHandlerStruct HealthHandlerInterface) KubeHandlerInterface {
	return &kubeHandlerStruct{ctx: ctx, healthHandlerStruct: healthHandlerStruct}
}

func (kh *kubeHandlerStruct) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Alive"})
}

func (kh *kubeHandlerStruct) GetReady(c *gin.Context) {
	// Kafka server check
	kafkaError := kh.healthHandlerStruct.GetKafkaHealth()
	if kafkaError != nil {
		logger.Errorf("Failed to dial kafka bootstrap server %v\n - Pod will stop serving request", kafkaError)
		c.JSON(http.StatusServiceUnavailable, gin.H{Status: "kafka bootstrap server health check failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{Status: Ready})
}

func (kh *HealthHandlerStruct) GetKafkaHealth() error {
	return NewKafkaHealthHandler().HandleHealth()
}
