// (C) Copyright 2022 Hewlett Packard Enterprise Development LP
package kubernetes

import (
	"context"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
)

const (
	TCP = "tcp"
)

var (
	logger = logging.GetLogger()
)

type Handler interface {
	HandleHealth() error
}

type handlerImpl struct {
	dialer *kafka.Dialer
}

func NewKafkaHealthHandler() Handler {
	return &handlerImpl{dialer: &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}}
}

func (hh *handlerImpl) HandleHealth() error {
	// Check whether required kafka server is reachable
	bootstrapServersSlice := strings.Split(configs.GetKafkaBootstrapServers(), `,`)
	var err error
	for i := 0; i < len(bootstrapServersSlice); i++ {
		err = hh.dialBootstrapServer(bootstrapServersSlice[i])
		if err == nil {
			// Atleast successful in creating topics on one bootstrap server
			logger.Debugf("Dial kafka bootstrap server %s is successful", bootstrapServersSlice[i])
			return nil
		}
	}
	return err
}

func (hh *handlerImpl) dialBootstrapServer(bootstrapServer string) error {
	conn, err := hh.dialer.DialContext(context.Background(), TCP, bootstrapServer)
	defer hh.closeKafkaConn(conn)
	if err != nil {
		logger.Errorf("Failed to dial kafka bootstrap server %s with err %v",
			bootstrapServer, err)
		return err
	}
	return nil
}

func (hh *handlerImpl) closeKafkaConn(conn *kafka.Conn) {
	if conn != nil {
		conn.Close()
	}
}
