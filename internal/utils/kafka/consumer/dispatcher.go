// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package consumer

import (
	"context"
	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/commonclient"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/commonclient/model"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/handlers"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/handlers/prometheus"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/services"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/kafka/producer"
	"go.uber.org/zap"
)

const (
	ErrorLabel      = "error"
	SuccessLabel    = "success"
	EndOfCollection = "EOC"
)

var (
	errEmptyID = errors.New("mandatory field is empty")
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
)

type Dispatcher func(
	context.Context,
	kafka.Message,
	producer.Producer,
	producer.Producer,
)

func validateSchedulerInput(ctx context.Context, input *model.KafkaCollectionRequest) error {
	if input.ApplicationCustomerID == "" || input.PlatformCustomerID == "" ||
		input.CollectionID == "" || input.CollectionType == "" || input.CollectionTrigger == "" ||
		input.Region == "" {
		logger.WithContext(ctx).Error("Fetched message validateSchedulerInput",
			zap.String("ApplicationCustomerID", input.ApplicationCustomerID),
			zap.String("PlatformCustomerID", input.PlatformCustomerID),
			zap.String("CollectionID", input.CollectionID),
			zap.String("CollectionType", input.CollectionType),
			zap.String("CollectionTrigger", input.CollectionTrigger),
			zap.String("Region", input.Region),
		)
		return errEmptyID
	}
	return nil
}

//nolint:funlen,dupl // required
func dispatchMessage(ctx context.Context, msg kafka.Message, schedulerProducer producer.Producer,
	harmonyProducer producer.Producer) {
	defer wg.Done()
	grpcContext, grpcCancel := context.WithCancel(ctx)
	defer grpcCancel()

	messageType := getMessageHeader(&msg, "ce_type")
	if messageType == "" {
		logger.WithContext(grpcContext).Error("Ignoring message with unspecified type")
		prometheus.KafkaConsumerEventsCnt.WithLabelValues(configs.GetHaulerConsumerKafkaTopic(), ErrorLabel).Inc()
		return
	}

	cID := getMessageHeader(&msg, "ce_collectionid")
	if cID == "" {
		logger.WithContext(grpcContext).Error("Ignoring message with missing ce_collectionid")
		prometheus.KafkaConsumerEventsCnt.WithLabelValues(configs.GetHaulerConsumerKafkaTopic(), ErrorLabel).Inc()
		return
	}

	data := &model.KafkaCollectionRequest{}
	err := json.Unmarshal(msg.Value, &data)
	if err != nil {
		logger.WithContext(grpcContext).Error(err.Error(), err)
		prometheus.KafkaConsumerEventsCnt.WithLabelValues(configs.GetHaulerConsumerKafkaTopic(), ErrorLabel).Inc()
		return
	}
	err = validateSchedulerInput(grpcContext, data)
	if err != nil {
		logger.WithContext(grpcContext).Error("Invalid input", zap.Error(err))
		prometheus.KafkaConsumerEventsCnt.WithLabelValues(configs.GetHaulerConsumerKafkaTopic(), ErrorLabel).Inc()
		return
	}

	validCollectionTypes := []string{"Inventory", "Volume-Perf"}
	if !handlers.Contains(validCollectionTypes, data.CollectionType) {
		logger.WithContext(grpcContext).Error("Invalid collection type specified")
		prometheus.KafkaConsumerEventsCnt.WithLabelValues(configs.GetHaulerConsumerKafkaTopic(), ErrorLabel).Inc()
		return
	}

	var consumerDetails = &handlers.ConsumerDetails{
		PlatformCustomerID:    data.PlatformCustomerID,
		CollectionID:          data.CollectionID,
		Region:                data.Region,
		ApplicationCustomerID: data.ApplicationCustomerID,
		ApplicationInstanceID: data.ApplicationInstanceID,
		CollectionTrigger:     data.CollectionTrigger,
		CollectionType:        data.CollectionType,
	}

	// metrics on total consumer events count and per customerID events count
	prometheus.KafkaConsumerEventsCnt.WithLabelValues(configs.GetHaulerConsumerKafkaTopic(), SuccessLabel).Inc()

	// var errorMap = make(map[string]map[string]string)
	// var statusMap = make(map[string]string)
	// var collectionStartTime = time.Now().UTC().String()
	// var errorMapEOC = make(map[string]map[string]string)

	arcusClient, arcusClientErr := commonclient.NewCommonClient(grpcContext)
	if arcusClient == nil {
		logger.WithContext(grpcContext).Errorf("Client initiation failed %v", arcusClientErr)
		// handlers.SetNested(errorMap, handlers.RESTError, handlers.RESTError, arcusClientErr.Error())
		// var data = handlers.ConstructFailedCollectionJSON(*consumerDetails, collectionStartTime,
		// 	handlers.UnknownDeviceType, handlers.ArcusHaulertype, errorMap)
		// file, _ := json.MarshalIndent(data, "", " ")
		// handlers.UploadToServer(grpcContext, file, *consumerDetails, services.Common, errorMap, schedulerProducer,
		// 	harmonyProducer)
		// consumerDetails.CollectionType = EndOfCollection
		// var eoc = handlers.ConstructEOCIndicatorJSON(*consumerDetails, collectionStartTime, handlers.ArcusHaulertype,
		// 	statusMap)
		// eocJSON, _ := json.MarshalIndent(eoc, "", " ")
		// handlers.UploadFileToServer(grpcContext, eocJSON, *consumerDetails, handlers.ArcusHaulertype, handlers.UnknownDeviceType,
		// 	statusMap, errorMapEOC, schedulerProducer, harmonyProducer)
		return
	}

	dataCollectionService := services.NewDataCollectionService(grpcContext, arcusClient)
	if dataCollectionService != nil {
		logger.WithContext(grpcContext).Infof("Collection started for collectionID: %s", consumerDetails.CollectionID)
		dataCollectionService.CollectDeviceInformation(grpcContext, *consumerDetails, schedulerProducer, harmonyProducer)
		logger.WithContext(grpcContext).Infof("Collection complete for collectionID: %s", consumerDetails.CollectionID)
	}
}
