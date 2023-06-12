// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/constants"
)

var KafkaDuration *prometheus.HistogramVec
var KafkaConsumerEventsCnt *prometheus.CounterVec
var KafkaProducerEventsCnt *prometheus.CounterVec
var UploadFileSize *prometheus.SummaryVec
var AwsRequestDuration *prometheus.HistogramVec

type MetricCollectionHandler interface {
	SetupMetricCollector()
}

type MetricCollectionHandlerImpl struct {
	kafkaDuration                *prometheus.HistogramVec
	kafkaConsumerEventsCnt       *prometheus.CounterVec
	kafkaProducerEventsCnt       *prometheus.CounterVec
	uploadFileSize               *prometheus.SummaryVec
	awsRequestProcessingDuration *prometheus.HistogramVec
}

func NewMetricCollectionHandler() MetricCollectionHandler {
	return &MetricCollectionHandlerImpl{}
}

func (mch *MetricCollectionHandlerImpl) SetupMetricCollector() {
	if KafkaDuration == nil {
		mch.kafkaDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: constants.PrometheusMetrics,
			Name:      "kafka_process_duration_seconds",
			Help:      "The total time in seconds of the kafka message processing operation.",
			Buckets:   []float64{5, 10, 20, 30, 40, 50},
		},
			[]string{"topic"})

		KafkaDuration = mch.kafkaDuration
	}

	if KafkaConsumerEventsCnt == nil {
		mch.kafkaConsumerEventsCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: constants.PrometheusMetrics,
			Name:      "kafka_consumer_events_count",
			Help:      "Number of total kafka consumer events received",
		},
			[]string{"topic", "status"})

		KafkaConsumerEventsCnt = mch.kafkaConsumerEventsCnt
	}

	if KafkaProducerEventsCnt == nil {
		mch.kafkaProducerEventsCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
			Subsystem: constants.PrometheusMetrics,
			Name:      "kafka_producer_events_count",
			Help:      "Number of total kafka producer events produced",
		},
			[]string{"topic", "deviceType", "collectionStatus", "uploadStatus", "collectionType"})

		KafkaProducerEventsCnt = mch.kafkaProducerEventsCnt
	}

	if UploadFileSize == nil {
		mch.uploadFileSize = prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Subsystem: constants.PrometheusMetrics,
			Name:      "upload_file_size_bytes",
			Help:      "Size of the file in bytes to be uploaded to upload server",
		},
			[]string{"collectionType", "deviceType"})

		UploadFileSize = mch.uploadFileSize
	}

	if AwsRequestDuration == nil {
		mch.awsRequestProcessingDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Subsystem: constants.PrometheusMetrics,
			Name:      "aws_upload_process_duration_seconds",
			Help:      "The total time in seconds to complete a Fleet Hauler upload (AWS API) request.",
			Buckets:   []float64{5, 10, 20, 30, 40, 50},
		},
			[]string{"deviceType", "collectionStatus", "uploadStatus"})

		AwsRequestDuration = mch.awsRequestProcessingDuration
	}
}
