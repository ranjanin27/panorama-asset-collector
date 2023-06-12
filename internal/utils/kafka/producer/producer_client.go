// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package producer

import (
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type WriterClient interface {
	WriteMessages(ctx context.Context, w *kafka.Writer, msgs []kafka.Message) error
	Close(r *kafka.Writer) error
}

type WriterClientImpl struct{}

func (pw *WriterClientImpl) WriteMessages(ctx context.Context, w *kafka.Writer, msgs []kafka.Message) error {
	return w.WriteMessages(ctx, msgs...)
}

func (pw *WriterClientImpl) Close(w *kafka.Writer) error {
	return w.Close()
}
