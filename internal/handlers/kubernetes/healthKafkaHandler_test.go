// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKafkaHealthHandler(t *testing.T) {
	handler := NewKafkaHealthHandler()
	err := handler.HandleHealth()
	assert.NotNil(t, err)
}
