// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/commonclient/model"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/handlers/prometheus"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/constants"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/kafka/producer"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
	"github.hpe.com/nimble-dcs/panorama-s3-upload/pkg/s3"
)

var (
	logger = logging.GetLogger()
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
)

const (
	Success                        = "Success"
	Partial                        = "Partial"
	Failed                         = "Failed"
	UnknownDeviceType              = "Unknown"
	Device1                        = "deviceType1"
	Device2                        = "deviceType2"
	HaulerType                     = "Fleet"
	ArcusHaulertype                = "Arcus"
	GRPCError                      = "GRPCError"
	RESTError                      = "RESTError"
	AuthError                      = "AuthError"
	Systems                        = "Systems"
	EndOfCollection                = "EOC"
	AuthzBearer             string = "Bearer"
	UserClaimName           string = "user_ctx"
	UserEmail               string = "sub"
	JWTParts                int    = 2
	tokenParts              int    = 3
	XForwardedFor                  = "x-forwarded-for"
	Bearer                         = "Bearer "
	maxNumPaddingBytesInJWT        = 4
)

func MapPFHErrorsToGRPCErrors(err error) (statusCode codes.Code, statusMsg string) {
	return MapPFHErrorsToGRPCErrorsWithContext(context.Background(), err)
}

//nolint:gocyclo // All error codes captured
func MapPFHErrorsToGRPCErrorsWithContext(ctx context.Context, err error) (statusCode codes.Code, statusMsg string) {
	if grpcErr, ok := status.FromError(err); ok {
		switch grpcErr.Code() {
		case codes.OK:
			return codes.OK, grpcErr.Message()
		case codes.Canceled:
			return codes.Canceled, grpcErr.Message()
		case codes.Unknown:
			return codes.Unknown, grpcErr.Message()
		case codes.InvalidArgument:
			return codes.InvalidArgument, grpcErr.Message()
		case codes.DeadlineExceeded:
			return codes.DeadlineExceeded, grpcErr.Message()
		case codes.NotFound:
			return codes.NotFound, grpcErr.Message()
		case codes.AlreadyExists:
			return codes.AlreadyExists, grpcErr.Message()
		case codes.PermissionDenied:
			return codes.PermissionDenied, grpcErr.Message()
		case codes.Unauthenticated:
			return codes.Unauthenticated, grpcErr.Message()
		case codes.ResourceExhausted:
			return codes.ResourceExhausted, grpcErr.Message()
		case codes.FailedPrecondition:
			return codes.FailedPrecondition, grpcErr.Message()
		case codes.Aborted:
			return codes.Aborted, grpcErr.Message()
		case codes.OutOfRange:
			return codes.OutOfRange, grpcErr.Message()
		case codes.Unimplemented:
			return codes.Unimplemented, grpcErr.Message()
		case codes.Internal:
			return codes.Internal, grpcErr.Message()
		case codes.Unavailable:
			return codes.Unavailable, grpcErr.Message()
		case codes.DataLoss:
			return codes.DataLoss, grpcErr.Message()
		}
	}
	logger.WithContext(ctx).Errorf("Unknown gRPC error: %v", err.Error())
	return codes.Internal, utils.InternalServerErrorMsg
}

func MapPFHErrorsToRESTErrors(err error) (statusCode int, statusMsg string) {
	return MapPFHErrorsToRESTErrorsWithContext(context.Background(), err)
}

func MapPFHErrorsToRESTErrorsWithContext(ctx context.Context, err error) (statusCode int, statusMsg string) {
	switch err := err.(type) {
	case utils.NotFoundError:
		return http.StatusNotFound, utils.NotFoundErrorMsg
	case utils.SoftDeleteError:
		return http.StatusNotFound, utils.NotFoundErrorMsg
	case utils.MissingAuthorizationError:
		return http.StatusBadRequest, utils.MissingAuthorizationErrorMsg
	case utils.InvalidAuthorizationError:
		return http.StatusBadRequest, err.Error()
	case utils.NotAuthorizedError: // This error is currently not used
		return http.StatusForbidden, utils.NotAuthorizedErrorMsg
	case utils.InvalidFilterError:
		return http.StatusBadRequest, err.Error()
	case utils.InvalidPathParamError:
		return http.StatusBadRequest, utils.InvalidPathParamErrorMsg
	case utils.BodyNotFoundError:
		return http.StatusBadRequest, utils.BodyNotFoundErrorMsg
	case utils.InvalidRequestFormatError:
		return http.StatusBadRequest, utils.InvalidRequestFormatErrorMsg
	case utils.RestClientError:
		return err.StatusCode, err.Error()
	default:
		// utils.InternalError and utils.DBConnError come under default
		// as we need to log in these cases and send internal server error
		logger.WithContext(ctx).Error(err.Error())
		return http.StatusInternalServerError, utils.InternalServerErrorMsg
	}
}

type ConsumerDetails struct {
	PlatformCustomerID, CollectionID, Region, ApplicationCustomerID, ApplicationInstanceID,
	CollectionTrigger, CollectionType string
}

type DeviceFailedCollection struct {
	Version, PlatformCustomerID, CollectionID, Region, ApplicationCustomerID, ApplicationInstanceID,
	CollectionTrigger, CollectionStartTime, CollectionEndTime, DeviceType, HaulerType, CollectionType string
	Error map[string]map[string]string
}

type EOCIndicator struct {
	Version, PlatformCustomerID, CollectionID, Region, ApplicationCustomerID, ApplicationInstanceID,
	CollectionTrigger, CollectionStartTime, CollectionEndTime, HaulerType, CollectionType string
	DeviceTypeUploadStatus map[string]string
}

type CommonHauler struct {
	Version, PlatformCustomerID, CollectionID, Region, ApplicationCustomerID, ApplicationInstanceID,
	CollectionTrigger, CollectionStartTime, CollectionEndTime, HaulerType, CollectionType string
	VirtualMachines, Datastores, ProtectionPolicies,
	VMProtectionGroups, ProtectedVMs, CSPMachineInstances []interface{}
	VMBackups   map[string][]interface{}
	VMSnapshots map[string][]interface{}
	DSBackups   map[string][]interface{}
	DSSnapshots map[string][]interface{}
	Error       map[string]map[string]string
}

type DeviceType2 struct {
	Version, PlatformCustomerID, CollectionID, Region, ApplicationCustomerID, ApplicationInstanceID,
	CollectionTrigger, CollectionStartTime, CollectionEndTime, DeviceType, HaulerType, CollectionType string
	DeviceTypeList    []string
	Systems           []interface{}
	Volumes           map[string][]interface{}
	Snapshots         map[string][]interface{}
	VolumeCollections map[string][]interface{}
	StoragePools      map[string][]interface{}
	VolumePerformance []interface{}
	Error             map[string]map[string]string
}

type DeviceType1 struct {
	Version, PlatformCustomerID, CollectionID, Region, ApplicationCustomerID, ApplicationInstanceID,
	CollectionTrigger, CollectionStartTime, CollectionEndTime, DeviceType, HaulerType, CollectionType string
	DeviceTypeList    []string
	Systems           []interface{}
	SystemCapacity    []interface{}
	Volumes           map[string][]interface{}
	Snapshots         map[string][]interface{}
	Vluns             map[string][]interface{}
	Applicationsets   map[string][]interface{}
	VolumePerformance []interface{}
	Error             map[string]map[string]string
}

type DeviceType4 struct {
	Version, PlatformCustomerID, CollectionID, Region, ApplicationCustomerID, ApplicationInstanceID,
	CollectionTrigger, CollectionStartTime, CollectionEndTime, DeviceType, HaulerType, CollectionType string
	DeviceTypeList    []string
	Systems           []interface{}
	SystemCapacity    []interface{}
	Volumes           map[string][]interface{}
	Clones            map[string][]interface{}
	Snapshots         map[string][]interface{}
	Vluns             map[string][]interface{}
	Applicationsets   map[string][]interface{}
	StoragePools      map[string][]interface{}
	VolumePerformance []interface{}
	Error             map[string]map[string]string
}

func ConstructCommonJSON(consumerDetails ConsumerDetails, collectionStartTime, haulerType string, virtualmachines,
	datatstores, protectpolicies, vmprotectgroups, protectedvms, cspmachineinstances []interface{},
	vmbackups, dsbackups map[string][]interface{}, errorMap map[string]map[string]string) CommonHauler {
	return CommonHauler{
		Version:               constants.Version,
		PlatformCustomerID:    consumerDetails.PlatformCustomerID,
		CollectionID:          consumerDetails.CollectionID,
		CollectionStartTime:   collectionStartTime,
		CollectionEndTime:     time.Now().UTC().String(),
		CollectionTrigger:     consumerDetails.CollectionTrigger,
		Region:                consumerDetails.Region,
		ApplicationCustomerID: consumerDetails.ApplicationCustomerID,
		ApplicationInstanceID: consumerDetails.ApplicationInstanceID,
		HaulerType:            haulerType,
		CollectionType:        consumerDetails.CollectionType,
		VirtualMachines:       virtualmachines,
		Datastores:            datatstores,
		ProtectionPolicies:    protectpolicies,
		VMProtectionGroups:    vmprotectgroups,
		ProtectedVMs:          protectedvms,
		CSPMachineInstances:   cspmachineinstances,
		VMBackups:             vmbackups,
		DSBackups:             dsbackups,
		Error:                 errorMap,
	}
}

func ConstructFailedCollectionJSON(consumerDetails ConsumerDetails, collectionStartTime, deviceType,
	haulerType string, errorMap map[string]map[string]string) DeviceFailedCollection {
	return DeviceFailedCollection{
		Version:               constants.Version,
		PlatformCustomerID:    consumerDetails.PlatformCustomerID,
		CollectionID:          consumerDetails.CollectionID,
		CollectionStartTime:   collectionStartTime,
		CollectionEndTime:     time.Now().UTC().String(),
		CollectionTrigger:     consumerDetails.CollectionTrigger,
		Region:                consumerDetails.Region,
		ApplicationCustomerID: consumerDetails.ApplicationCustomerID,
		ApplicationInstanceID: consumerDetails.ApplicationInstanceID,
		HaulerType:            haulerType,
		DeviceType:            deviceType,
		CollectionType:        consumerDetails.CollectionType,
		Error:                 errorMap,
	}
}

func ConstructEOCIndicatorJSON(consumerDetails ConsumerDetails, collectionStartTime, haulerType string,
	statusMap map[string]string) EOCIndicator {
	return EOCIndicator{
		Version:                constants.Version,
		PlatformCustomerID:     consumerDetails.PlatformCustomerID,
		CollectionID:           consumerDetails.CollectionID,
		CollectionStartTime:    collectionStartTime,
		CollectionEndTime:      time.Now().UTC().String(),
		CollectionTrigger:      consumerDetails.CollectionTrigger,
		Region:                 consumerDetails.Region,
		ApplicationCustomerID:  consumerDetails.ApplicationCustomerID,
		ApplicationInstanceID:  consumerDetails.ApplicationInstanceID,
		HaulerType:             haulerType,
		CollectionType:         consumerDetails.CollectionType,
		DeviceTypeUploadStatus: statusMap,
	}
}

func ConstructDeviceType2JSON(consumerDetails ConsumerDetails, collectionStartTime, deviceType,
	haulerType string, deviceTypeList []string, systems, volumePerformance []interface{}, volumes,
	volumeCollections, storagePools, snapshots map[string][]interface{},
	errorMap map[string]map[string]string) DeviceType2 {
	return DeviceType2{
		Version:               constants.Version,
		PlatformCustomerID:    consumerDetails.PlatformCustomerID,
		CollectionID:          consumerDetails.CollectionID,
		CollectionStartTime:   collectionStartTime,
		CollectionEndTime:     time.Now().UTC().String(),
		CollectionTrigger:     consumerDetails.CollectionTrigger,
		Region:                consumerDetails.Region,
		ApplicationCustomerID: consumerDetails.ApplicationCustomerID,
		ApplicationInstanceID: consumerDetails.ApplicationInstanceID,
		HaulerType:            haulerType,
		DeviceType:            deviceType,
		CollectionType:        consumerDetails.CollectionType,
		DeviceTypeList:        deviceTypeList,
		Systems:               systems,
		Volumes:               volumes,
		Snapshots:             snapshots,
		VolumeCollections:     volumeCollections,
		StoragePools:          storagePools,
		VolumePerformance:     volumePerformance,
		Error:                 errorMap,
	}
}

func ConstructDeviceType1JSON(consumerDetails ConsumerDetails, collectionStartTime, deviceType,
	haulerType string, deviceTypeList []string, systems, systemCapacity, volumePerformance []interface{},
	volumes, applicationsets, snapshots, vluns map[string][]interface{},
	errorMap map[string]map[string]string) DeviceType1 {
	return DeviceType1{
		Version:               constants.Version,
		PlatformCustomerID:    consumerDetails.PlatformCustomerID,
		CollectionID:          consumerDetails.CollectionID,
		CollectionStartTime:   collectionStartTime,
		CollectionEndTime:     time.Now().UTC().String(),
		CollectionTrigger:     consumerDetails.CollectionTrigger,
		Region:                consumerDetails.Region,
		ApplicationCustomerID: consumerDetails.ApplicationCustomerID,
		ApplicationInstanceID: consumerDetails.ApplicationInstanceID,
		HaulerType:            haulerType,
		DeviceType:            deviceType,
		CollectionType:        consumerDetails.CollectionType,
		DeviceTypeList:        deviceTypeList,
		Systems:               systems,
		SystemCapacity:        systemCapacity,
		Volumes:               volumes,
		Snapshots:             snapshots,
		Applicationsets:       applicationsets,
		Vluns:                 vluns,
		VolumePerformance:     volumePerformance,
		Error:                 errorMap,
	}
}

func ConstructDeviceType4JSON(consumerDetails ConsumerDetails, collectionStartTime, deviceType,
	haulerType string, deviceTypeList []string, systems, systemCapacity, volumePerformance []interface{},
	volumes, applicationsets, snapshots, vluns, clones, storagePools map[string][]interface{},
	errorMap map[string]map[string]string) DeviceType4 {
	return DeviceType4{
		Version:               constants.Version,
		PlatformCustomerID:    consumerDetails.PlatformCustomerID,
		CollectionID:          consumerDetails.CollectionID,
		CollectionStartTime:   collectionStartTime,
		CollectionEndTime:     time.Now().UTC().String(),
		CollectionTrigger:     consumerDetails.CollectionTrigger,
		Region:                consumerDetails.Region,
		ApplicationCustomerID: consumerDetails.ApplicationCustomerID,
		ApplicationInstanceID: consumerDetails.ApplicationInstanceID,
		HaulerType:            haulerType,
		DeviceType:            deviceType,
		CollectionType:        consumerDetails.CollectionType,
		DeviceTypeList:        deviceTypeList,
		Systems:               systems,
		SystemCapacity:        systemCapacity,
		Volumes:               volumes,
		Clones:                clones,
		Snapshots:             snapshots,
		Applicationsets:       applicationsets,
		Vluns:                 vluns,
		StoragePools:          storagePools,
		VolumePerformance:     volumePerformance,
		Error:                 errorMap,
	}
}

func SetNested(m1 map[string]map[string]string, k1, k2, v string) {
	m2 := m1[k1]
	if m2 == nil {
		m2 = make(map[string]string)
		m1[k1] = m2
	}
	m2[k2] = v
}

func ClearMap(m map[string][]interface{}) {
	for k := range m {
		delete(m, k)
	}
}

func ClearErrorMap(m map[string]map[string]string) {
	for k := range m {
		delete(m, k)
	}
}

func ClearMainMapEntries(mainVolumeCollectionMap2, mainStoragePoolMap2, mainVolumeMap2, mainVolumeSnapshotsMap2, mainVolumeMap1,
	mainApplicationSetMap1, mainVolumeSnapshotsMap1, mainVlunMap1, mainApplicationSetMap4, mainCloneMap4, mainStoragePoolMap4,
	mainVlunMap4, mainVolumeMap4, mainVolumeSnapshotsMap4 map[string][]interface{}, mainErrorMap2,
	mainErrorMap1, mainErrorMap4 map[string]map[string]string) {
	ClearMap(mainVolumeCollectionMap2)
	ClearMap(mainStoragePoolMap2)
	ClearMap(mainVolumeMap2)
	ClearMap(mainVolumeSnapshotsMap2)
	ClearMap(mainVolumeMap1)
	ClearMap(mainApplicationSetMap1)
	ClearMap(mainVolumeSnapshotsMap1)
	ClearMap(mainVlunMap1)
	ClearErrorMap(mainErrorMap1)
	ClearErrorMap(mainErrorMap2)
}

func PublishtoHarmonyKafka(ctx context.Context, consumerDetails ConsumerDetails, awsKey, bucketName string,
	fileSize int, harmonyProducer producer.Producer) {
	res := model.HarmonyStatusRequest{}
	res.Source.S3.Bucket = s3.GetS3BucketFromURL(bucketName)
	res.Source.S3.Key = awsKey
	res.Notification.FileDomain = configs.GetFileDomain()
	res.Notification.FileType = configs.GetSourceType()
	res.Notification.FileSize = fileSize
	res.Notification.EntityID = consumerDetails.CollectionID + "_" + time.Now().UTC().Format(time.RFC3339)
	res.Notification.FileName = awsKey
	encodedResponse, err := json.Marshal(res)
	if err != nil {
		logger.WithContext(ctx).Warnf("Could not marshal the input to JSON for topic %v",
			configs.GetHaulerHarmonyKafkaTopic())
	}
	err = harmonyProducer.PublishMessage(ctx, consumerDetails.CollectionID, consumerDetails.ApplicationCustomerID,
		configs.GetKafkaEventType(), configs.GetKafkaEventSource(), time.Now().UTC(), encodedResponse)
	if err != nil {
		logger.WithContext(ctx).Errorf("Publishing kafka message failed from hauler to topic %v",
			configs.GetHaulerHarmonyKafkaTopic())
	} else {
		logger.WithContext(ctx).Infof("Successfully published to topic %v",
			configs.GetHaulerHarmonyKafkaTopic())
		prometheus.KafkaProducerEventsCnt.WithLabelValues(configs.GetHaulerHarmonyKafkaTopic(),
			UnknownDeviceType, Success, Success, EndOfCollection).Inc()
	}
}

func PublishtoCommonKafka(ctx context.Context, consumerDetails ConsumerDetails, awsKey, bucketName string,
	fileSize int, harmonyProducer producer.Producer) {
	res := model.HarmonyStatusRequest{}
	res.Source.S3.Bucket = s3.GetS3BucketFromURL(bucketName)
	res.Source.S3.Key = awsKey
	res.Notification.FileDomain = configs.GetFileDomain()
	res.Notification.FileType = "Common"
	res.Notification.FileSize = fileSize
	res.Notification.EntityID = consumerDetails.CollectionID + "_" + time.Now().UTC().Format(time.RFC3339)
	res.Notification.FileName = awsKey
	encodedResponse, err := json.Marshal(res)
	if err != nil {
		logger.WithContext(ctx).Warnf("Could not marshal the input to JSON for topic %v",
			configs.GetHaulerHarmonyKafkaTopic())
	}
	err = harmonyProducer.PublishMessage(ctx, consumerDetails.CollectionID, consumerDetails.ApplicationCustomerID,
		configs.GetKafkaEventType(), configs.GetKafkaEventSource(), time.Now().UTC(), encodedResponse)
	if err != nil {
		logger.WithContext(ctx).Errorf("Publishing kafka message failed from hauler to topic %v",
			configs.GetHaulerHarmonyKafkaTopic())
	} else {
		logger.WithContext(ctx).Infof("Successfully published to topic %v", configs.GetHaulerHarmonyKafkaTopic())
	}
}

func PublishtoSchedulerKafka(ctx context.Context, consumerDetails ConsumerDetails, haulerType,
	uploadStatus, collectionStatus, deviceType, fileSize, keyName, errors string, schedulerProducer producer.Producer) {
	res := &model.ScheduleStatus{
		CollectionID:          consumerDetails.CollectionID,
		PlatformCustomerID:    consumerDetails.PlatformCustomerID,
		ApplicationCustomerID: consumerDetails.ApplicationCustomerID,
		CollectionType:        consumerDetails.CollectionType,
		HaulerType:            haulerType,
		DeviceType:            deviceType,
		JSONVersion:           constants.Version,
		UploadStatus:          uploadStatus,
		CollectionStatus:      collectionStatus,
		UploadFileSize:        fileSize,
		S3Bucket:              configs.GetAWSS3BucketName(),
		FileName:              keyName,
		Error:                 errors,
	}
	encodedResponse, err := json.Marshal(res)
	if err != nil {
		logger.WithContext(ctx).Warnf("Could not marshal the input to JSON for topic %v",
			configs.GetHaulerProducerKafkaTopic())
	}

	err = schedulerProducer.PublishMessage(ctx, consumerDetails.CollectionID, consumerDetails.ApplicationCustomerID,
		configs.GetKafkaEventType(), configs.GetKafkaEventSource(), time.Now().UTC(), encodedResponse)
	if err != nil {
		logger.WithContext(ctx).Errorf("Publishing kafka message failed from hauler to topic %v",
			configs.GetHaulerProducerKafkaTopic())
	} else {
		logger.WithContext(ctx).Infof("Successfully published to topic %v",
			configs.GetHaulerProducerKafkaTopic())
		prometheus.KafkaProducerEventsCnt.WithLabelValues(configs.GetHaulerProducerKafkaTopic(),
			deviceType, collectionStatus, uploadStatus, consumerDetails.CollectionType).Inc()
	}
}

var UploadToS3 = func(jsonContent *[]byte, bucketName, awsS3Region, awsAccessKeyID, awsSecretAccessKey string,
	uploadType s3.UploadType) (string, int, error) {
	keyname, filesize, err := s3.UploadToAwsS3(jsonContent, bucketName, awsS3Region, awsAccessKeyID,
		awsSecretAccessKey, s3.UploadType(configs.GetSourceType()))
	return keyname, filesize, err
}

func UploadFileToServer(ctx context.Context, jsonContent []byte, consumerDetails ConsumerDetails, haulerType,
	deviceType string, statusMap map[string]string, errorMap map[string]map[string]string,
	schedulerProducer producer.Producer, harmonyProducer producer.Producer) {
	var errors string
	var uploadStatus = Success
	var collectionStatus string
	if len(errorMap) == 0 {
		collectionStatus = Success
		errors = ""
	} else if v, found := errorMap[GRPCError]; found {
		collectionStatus = Failed
		errors = v[GRPCError]
	} else if v, found := errorMap[RESTError]; found {
		collectionStatus = Failed
		errors = v[RESTError]
	} else if v, found := errorMap[AuthError]; found {
		collectionStatus = Failed
		errors = v[AuthError]
	} else if _, found := errorMap[Systems]; found {
		collectionStatus = Failed
		bs, _ := json.Marshal(errorMap)
		errors = string(bs)
	} else {
		collectionStatus = Partial
		bs, _ := json.Marshal(errorMap)
		errors = string(bs)
	}

	requestStartTime := time.Now()
	keyName, filesize, err := UploadToS3(&jsonContent, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
		configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	if err != nil {
		logger.WithContext(ctx).Errorf("File upload failed with error %v", err)
		uploadStatus = Failed
	}
	logger.WithContext(ctx).Infof("CollectionType %v  CollectionID  %v CollectionStatus %v UploadStatus %v",
		consumerDetails.CollectionType, consumerDetails.CollectionID, collectionStatus, uploadStatus)
	elapsed := time.Since(requestStartTime).Seconds()
	prometheus.AwsRequestDuration.WithLabelValues(deviceType, collectionStatus, uploadStatus).Observe(elapsed)
	if uploadStatus != Failed {
		PublishtoHarmonyKafka(ctx, consumerDetails, keyName, configs.GetAWSS3BucketName(), filesize, harmonyProducer)
		prometheus.UploadFileSize.WithLabelValues(consumerDetails.CollectionType, deviceType).Observe(float64(filesize))
	}
	statusMap[deviceType] = uploadStatus
	PublishtoSchedulerKafka(ctx, consumerDetails, haulerType, uploadStatus, collectionStatus, deviceType,
		strconv.Itoa(filesize), keyName, errors, schedulerProducer)
}

func UploadToServer(ctx context.Context, jsonContent []byte, consumerDetails ConsumerDetails, haulerType string, errorMap map[string]map[string]string,
	schedulerProducer producer.Producer, harmonyProducer producer.Producer) {
	var errors string
	var uploadStatus = Success
	var collectionStatus string
	if len(errorMap) == 0 {
		collectionStatus = Success
		errors = ""
	} else if v, found := errorMap[RESTError]; found {
		collectionStatus = Failed
		errors = v[RESTError]
	} else if v, found := errorMap[AuthError]; found {
		collectionStatus = Failed
		errors = v[AuthError]
	} else {
		collectionStatus = Partial
		bs, _ := json.Marshal(errorMap)
		errors = string(bs)
	}

	keyName, filesize, err := UploadToS3(&jsonContent, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
		configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	if err != nil {
		logger.WithContext(ctx).Errorf("File upload failed with error %v", err)
		uploadStatus = Failed
	}
	logger.WithContext(ctx).Infof("CollectionType %v  CollectionID  %v CollectionStatus %v UploadStatus %v Error %v",
		consumerDetails.CollectionType, consumerDetails.CollectionID, collectionStatus, uploadStatus, errors)
	if uploadStatus != Failed {
		PublishtoCommonKafka(ctx, consumerDetails, keyName, configs.GetAWSS3BucketName(), filesize, harmonyProducer)
	}
	// PublishtoSchedulerKafka(ctx, consumerDetails, haulerType, uploadStatus, collectionStatus, deviceType,
	// 	strconv.Itoa(filesize), keyName, errors, schedulerProducer)
}

func Unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GetAuthorizationHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", utils.MissingAuthorizationError(-1)
	}
	return authHeader, nil
}
