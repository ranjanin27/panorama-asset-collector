// Copyright (c) 2022 Hewlett Packard Enterprise Development LP

package consumer

import (
	"context"
	"strings"
	"testing"

	kafka "github.com/segmentio/kafka-go"
)

func TestReaderClientImpl_NewReader(t *testing.T) {
	test := ReaderClientImpl{}
	brokers := "test-broker-1,test-broker-2"
	cfg := &kafka.ReaderConfig{
		Brokers:     strings.Split(brokers, ","),
		GroupTopics: []string{"some"},
		GroupID:     "groupID",
		StartOffset: kafka.FirstOffset,
	}
	test.NewReader(cfg)
}

func TestReaderClientImpl_FetchMessage(t *testing.T) {
	test := ReaderClientImpl{}
	brokers := "test-broker-3,test-broker-4"
	cfg := &kafka.ReaderConfig{
		Brokers:     strings.Split(brokers, ","),
		GroupTopics: []string{"some"},
		GroupID:     "groupID",
		StartOffset: kafka.FirstOffset,
	}
	ctx := context.Background()
	reader := test.NewReader(cfg)
	reader.Close()
	_, _ = test.FetchMessage(ctx, reader)
}

func TestReaderClientImpl_Close(t *testing.T) {
	test := ReaderClientImpl{}
	brokers := "test-broker-5,test-broker-6"
	cfg := &kafka.ReaderConfig{
		Brokers:     strings.Split(brokers, ","),
		GroupTopics: []string{"some"},
		GroupID:     "groupID",
		StartOffset: kafka.FirstOffset,
	}
	reader := test.NewReader(cfg)
	test.Close(reader)
}
