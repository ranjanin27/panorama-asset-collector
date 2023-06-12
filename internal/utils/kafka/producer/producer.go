// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package producer

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/contextutilities"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/tracer"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
	"github.hpe.com/cloud/go-gadgets/x/uuidgenerator"
	tildeerrors "github.hpe.com/cloud/tilde-common/pkg/errors"
	"go.uber.org/zap"
)

var (
	logger = logging.GetLogger()
)

const (
	messageMimeType = "application/json"
)

type Producer interface {
	CloseWriter() error

	PublishMessage(
		ctx context.Context,
		key string,
		customerID string,
		eventType string,
		eventSource string,
		eventTime time.Time,
		data []byte,
	) error
}

type Impl struct {
	client  WriterClient
	writer  *kafka.Writer
	uuidGen uuidgenerator.Generator
}

func (p *Impl) CloseWriter() error {
	logger.Info("Writer close")
	return p.client.Close(p.writer)
}

func NewProducer(sslMode bool, brokers, topic string, uuidGen uuidgenerator.Generator) Producer {
	writer := &kafka.Writer{
		Topic:        topic,
		Addr:         kafka.TCP(strings.Split(brokers, ",")[0]),
		RequiredAcks: kafka.RequireAll,
	}
	// This will be false during development
	if sslMode {
		writer.Transport = &kafka.Transport{TLS: &tls.Config{MinVersion: tls.VersionTLS12}}
	}
	return &Impl{
		client:  &WriterClientImpl{},
		writer:  writer,
		uuidGen: uuidGen,
	}
}

func (p *Impl) PublishMessage(
	ctx context.Context,
	collectionID, customerID, eventType, eventSource string, eventTime time.Time,
	data []byte) error {
	id, err := p.uuidGen.New()
	if err != nil {
		logger.WithContext(ctx).Error("Could not generate UUID", zap.Error(err))
		return tildeerrors.NewInternalError("unable to generate UUID", err)
	}
	traceProperties := contextutilities.ExtractTraceDetails(ctx)
	traceParent := tracer.MakeCETraceParent((*tracer.TraceHeaders)(traceProperties))
	traceState := tracer.MakeCETraceState((*tracer.TraceHeaders)(traceProperties))

	msgs := []kafka.Message{{
		Headers: []protocol.Header{
			{Key: "ce_specversion", Value: []byte("1.0")},
			{Key: "ce_type", Value: []byte(eventType)},
			{Key: "ce_id", Value: []byte(id.String())},
			{Key: "ce_source", Value: []byte(eventSource)},
			{Key: "ce_time", Value: []byte(eventTime.UTC().Format(time.RFC3339))},
			{Key: "content-type", Value: []byte(messageMimeType)},
			{Key: "ce_collectionid", Value: []byte(collectionID)},
			{Key: "ce_customerid", Value: []byte(customerID)},
			{Key: "ce_traceparent", Value: []byte(traceParent)},
			{Key: "ce_tracestate", Value: []byte(traceState)},
		},
		Key: []byte(customerID),
		// The value is expected to be a cloud event encoded in JSON.
		Value: data,
	}}

	err = p.client.WriteMessages(ctx, p.writer, msgs)
	if err != nil {
		logger.WithContext(ctx).Error("publishing message failed ", zap.Error(err))
		return tildeerrors.NewInternalError("Kafka write failed", err)
	}

	logger.WithContext(ctx).Info("Message written",
		zap.String("key", string(msgs[0].Key)),
		zap.String("value", string(msgs[0].Value)))
	return nil
}
