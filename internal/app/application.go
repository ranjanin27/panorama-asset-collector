// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package app

import (
	"context"
	"net/http"

	// This blank import is for enabling pprof
	// This should be removed in production.
	//nolint:gosec // G108, Adding this now for debugging purpose
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.hpe.com/cloud/go-gadgets/x/uuidgenerator"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/driver"
	metrics "github.hpe.com/nimble-dcs/panorama-common-hauler/internal/handlers/prometheus"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/kafka/consumer"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/kafka/producer"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
)

var (
	// Gin is a HTTP Web Framework in Go
	// https://github.com/gin-gonic/gin
	// Create a Gin router
	router = gin.New()
	logger = logging.GetLogger()
	cancel = func() {
		// Nothing to cancel.
	}
)

// StartRESTService will initialize and start the service
func startRESTService(ctx context.Context) {
	mapRESTUrls(ctx)
	if err := router.Run(":" + configs.GetHTTPPort()); err != nil {
		logger.Error("Failed to start router due to err...", err)
		return
	}
}

func CloseKafkaProducers(kafkaProducer producer.Producer) {
	if err := kafkaProducer.CloseWriter(); err != nil {
		logger.Info("failed to close producers:", err.Error())
	}
}

func SetGinLogMode() {
	// Skip logging the health service url(s) to avoid flood in the log messages
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter,
			"/healthz/liveness",
			"/healthz/readiness"),
		gin.Recovery(),
	)
	gin.SetMode(configs.GetGinLogLevel())
}

//nolint:gosec // G114
func enablePprof() {
	// Enable pprof in port 6060
	err := http.ListenAndServe("localhost:6060", nil)
	if err != nil {
		logger.Errorf("Error while starting pprof: %+v", err)
	}
}

func handleShutdownSignal() {
	logger.Info("Shutdown signal received")
	cancel()
}

func StartApplication(version string) {
	logger.Info("Starting application now...")
	SetGinLogMode()
	defer logger.Close()

	logger.Infof("Panorama Common Hauler service Version=%s starting...", version)

	if configs.GetPprofEnabled() {
		go enablePprof()
		logger.Infof("Pprof is enabled")
	}

	ctx := context.Background()
	ctx, cancel = context.WithCancel(ctx)
	driver.RegisterShutdownSignalHandler(handleShutdownSignal)

	logger.Info("Starting REST Server")
	go startRESTService(ctx)

	metricCollectionHandler := metrics.NewMetricCollectionHandler()
	metricCollectionHandler.SetupMetricCollector()
	initMetrics()

	kafkaSchedulerProducer := producer.NewProducer(configs.GetKafkaTLSOn(),
		configs.GetKafkaBootstrapServers(), configs.GetHaulerProducerKafkaTopic(),
		uuidgenerator.NewGoogleUUIDGenerator())

	kafkaHarmonyProducer := producer.NewProducer(configs.GetKafkaTLSOn(),
		configs.GetKafkaBootstrapServers(), configs.GetHaulerHarmonyKafkaTopic(),
		uuidgenerator.NewGoogleUUIDGenerator())

	kafkaConsumerDone := make(chan bool, 1)
	kafkaConsumer := consumer.NewConsumer(
		configs.GetHaulerConsumerKafkaTopic(),
		configs.GetHaulerKafkaGroupID(),
		configs.GetKafkaTLSOn(),
		configs.GetKafkaBootstrapServers(),
		kafkaSchedulerProducer,
		kafkaHarmonyProducer,
	)
	go kafkaConsumer.ReadAndProcessMessages(ctx, kafkaConsumerDone)

	// Block until a shutdown signal is received.
	<-ctx.Done()
	CloseKafkaProducers(kafkaSchedulerProducer)
	CloseKafkaProducers(kafkaHarmonyProducer)
	// Allow the Kafka consumer to clean up before exiting
	<-kafkaConsumerDone
}

func initMetrics() {
	prometheus.MustRegister(metrics.KafkaDuration)
	prometheus.MustRegister(metrics.KafkaConsumerEventsCnt)
	prometheus.MustRegister(metrics.KafkaProducerEventsCnt)
	prometheus.MustRegister(metrics.UploadFileSize)
	prometheus.MustRegister(metrics.AwsRequestDuration)
}
