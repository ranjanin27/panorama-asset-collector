// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package consumer

import (
	"context"
	"crypto/tls"
	"strings"
	"sync"
	"time"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/kafka/producer"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/tracer"
	"google.golang.org/grpc/metadata"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var (
	logger                       = logging.GetLogger()
	wg                           sync.WaitGroup
	msgCount                     = 0
	msgMap                       = make(map[int]kafka.Message)
	schedulerMessageToDispatcher = Dispatcher(dispatchMessage)
)

const (
	MinimumCommitCount = 5
)

type Consumer interface {
	ReadAndProcessMessages(ctx context.Context, done chan bool)
}

type Impl struct {
	topic             string
	groupID           string
	sslMode           bool
	brokers           string
	client            ReaderClient
	schedulerProducer producer.Producer
	harmonyProducer   producer.Producer
}

// NewConsumer creates a new Kafka consumer for a specified topic
func NewConsumer(
	topic string,
	groupID string,
	sslMode bool,
	brokers string,
	schedulerProducer producer.Producer,
	harmonyProducer producer.Producer,
) Consumer {
	return &Impl{
		topic:             topic,
		groupID:           groupID,
		sslMode:           sslMode,
		brokers:           brokers,
		client:            &ReaderClientImpl{},
		schedulerProducer: schedulerProducer,
		harmonyProducer:   harmonyProducer,
	}
}

func buildReaderConfig(topic, groupID string, sslMode bool, brokers string) *kafka.ReaderConfig {
	cfg := &kafka.ReaderConfig{
		Brokers:     strings.Split(brokers, ","),
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,    // Get message as soon as it is available
		MaxBytes:    1000, // Batch size
	}
	if sslMode {
		cfg.Dialer = &kafka.Dialer{
			TLS:       &tls.Config{MinVersion: tls.VersionTLS12},
			Timeout:   time.Duration(configs.GetKafkaTimeoutInSeconds()) * time.Second,
			DualStack: configs.GetKafkaDualStackOn(),
		}
	}
	return cfg
}

func (c *Impl) ReadAndProcessMessages(
	ctx context.Context, done chan bool) {
	configuration := buildReaderConfig(c.topic, c.groupID, c.sslMode, c.brokers)
	logger.WithContext(ctx).Info("Reading from kafka",
		zap.String("groupID", c.groupID),
		zap.Bool("sslMode", c.sslMode),
		zap.String("brokers", c.brokers),
	)
	reader := c.client.NewReader(configuration)
	logger.WithContext(ctx).Info("Kafka message reader is created")
	defer c.client.Close(reader)

	for ctx.Err() == nil {
		c.readAndDispatchMessage(ctx, reader, c.schedulerProducer, c.harmonyProducer)
	}
	logger.WithContext(ctx).Info("Kafka message reader is exiting")
	done <- true
}

func (c *Impl) readAndDispatchMessage(ctx context.Context,
	reader *kafka.Reader, schedulerProducer producer.Producer,
	harmonyProducer producer.Producer) {
	msg, err := c.client.FetchMessage(ctx, reader)
	if err != nil {
		logger.WithContext(ctx).Error("Failed to fetch message", zap.Error(err))
		return
	}
	traceParent := getMessageHeader(&msg, "ce_traceparent")
	traceState := getMessageHeader(&msg, "ce_tracestate")
	eventData := make(map[string]string)
	eventData[tracer.CETraceParent] = traceParent
	eventData[tracer.CETraceState] = traceState
	hdrs := tracer.ExtractFromIncomingKafkaEvent(eventData, logger)
	ctx = metadata.NewIncomingContext(
		ctx,
		metadata.Pairs(tracer.GRPCTraceID, hdrs.TraceID, tracer.GRPCSpanID, hdrs.SpanID,
			tracer.GRPCParentSpanID, hdrs.ParentSpanID, tracer.GRPCRequestID, hdrs.RequestID,
			tracer.GRPCSampled, "true"),
	)
	logger.WithContext(ctx).Info("Fetched message",
		zap.String("topic", c.topic),
		zap.Int("partition", msg.Partition),
		zap.Int64("offset", msg.Offset),
		zap.String("key", string(msg.Key)),
		zap.String("value", string(msg.Value)))

	// dispatch message here to hauler for collection
	dispatcher := schedulerMessageToDispatcher
	wg.Add(1)
	go dispatcher(ctx, msg, schedulerProducer, harmonyProducer)

	msgMap[msg.Partition] = msg
	msgCount++

	if msgCount%MinimumCommitCount == 0 {
		wg.Wait()
		for key := range msgMap {
			commitMsg := msgMap[key]
			err = c.client.CommitMessage(ctx, reader, &commitMsg)
			if err != nil {
				logger.WithContext(ctx).Error("Failed to commit message", zap.Error(err))
				return
			}
			delete(msgMap, key)
			logger.WithContext(ctx).Info("Committed message",
				zap.Int("partition", commitMsg.Partition),
				zap.Int64("offset", commitMsg.Offset),
				zap.String("key", string(commitMsg.Key)),
				zap.String("value", string(commitMsg.Value)),
			)
		}
	}
}

func getMessageHeader(message *kafka.Message, key string) string {
	for _, header := range message.Headers {
		if header.Key == key {
			return string(header.Value)
		}
	}
	return ""
}
