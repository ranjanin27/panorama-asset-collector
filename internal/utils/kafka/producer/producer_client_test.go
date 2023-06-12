// Copyright (c) 2022 Hewlett Packard Enterprise Development LP

package producer

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestWriterClientImpl_WriteMessages(t *testing.T) {
	ctx := context.Background()
	writer := WriterClientImpl{}
	_ = writer.WriteMessages(ctx, &kafka.Writer{}, []kafka.Message{})
}

func TestWriterClientImpl_Close(t *testing.T) {
	writer := WriterClientImpl{}
	writer.Close(&kafka.Writer{})
}
