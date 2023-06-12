// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package consumer

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// ReaderClient provides a mockable interface for kafka functions
// required to implement the Consumer interface.
type ReaderClient interface {
	NewReader(config *kafka.ReaderConfig) *kafka.Reader
	FetchMessage(ctx context.Context, r *kafka.Reader) (kafka.Message, error)
	CommitMessage(ctx context.Context, r *kafka.Reader, msg *kafka.Message) error
	Close(r *kafka.Reader) error
}

type ReaderClientImpl struct{}

func (cw *ReaderClientImpl) NewReader(config *kafka.ReaderConfig) *kafka.Reader {
	return kafka.NewReader(*config)
}

func (cw *ReaderClientImpl) FetchMessage(ctx context.Context, r *kafka.Reader) (kafka.Message, error) {
	return r.FetchMessage(ctx)
}

func (cw *ReaderClientImpl) CommitMessage(ctx context.Context, r *kafka.Reader, msg *kafka.Message) error {
	return r.CommitMessages(ctx, *msg)
}

func (cw *ReaderClientImpl) Close(r *kafka.Reader) error {
	return r.Close()
}
