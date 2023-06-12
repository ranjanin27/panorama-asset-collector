// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package configs

import (
	"time"

	"github.com/spf13/viper"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/constants"
)

const (
	/*
	* ===== Important: Viper configuration keys are case insensitive =====
	 */

	// Config Key Names
	// panorama-common-hauler
	httpPort    = "http_port"
	logLevel    = "log_level"
	ginlogLevel = "gin_log_level"

	config       = "CONFIG"
	debugMode    = "DEBUG_MODE"
	deployMode   = "DEPLOY_MODE" // test | prod
	pprofEnabled = "PPROF_ENABLED"
	localCluster = "local_cluster"

	// grpc
	grpcConnectionTimeout = "GRPC_TIMEOUT"

	// Authz configs
	authzEndpoint                      = "authz_endpoint"
	authzDisable                       = "authz_disable"
	updateServiceAuthExpiryTimeoutMins = "update_service_auth_expiry_timeout_mins"

	// Grip mock timeout
	gripmockConnectionTimeout = "gripmock_connection_timeout"

	// Customer ID
	useCustomerID = "use_customer_id"

	// kafka
	HaulerKafkaGroupID       = "hauler_kafka_group_id"
	HaulerConsumerKafkaTopic = "hauler_consumer_kafka_topic_name"
	HaulerProducerKafkaTopic = "hauler_producer_kafka_topic_name"
	HaulerHarmonyKafkaTopic  = "hauler_harmony_kafka_topic_name"
	kafkaBootstrapServers    = "kafka_bootstrap_servers"
	kafkaDualStack           = "kafka_dual_stack"
	kafkaBatchSize           = "kafka_batch_size"
	kafkaTLSOn               = "kafka_tls_mode"
	kafkaTimeoutInSeconds    = "kafka_timeout_in_seconds"

	kafkaEventSource = "kafka_event_source"
	kafkaEventType   = "kafka_event_type"

	// Fleet Grpc Endpoints
	fleetGrpcEndpoint            = "fleet_grpc_endpoint"
	fleetGrpcDeviceType1Endpoint = "fleet_grpc_devicetype1_endpoint"
	fleetGrpcDeviceType2Endpoint = "fleet_grpc_devicetype2_endpoint"
	// Internal Jwt
	internalJwt = "internal_jwt"
	// AWS upload
	AWSAccessKey     = "aws_access_key_id"
	AWSScrtAccessKey = "aws_secret_access_key"
	AWSRegion        = "aws_region"
	AWSS3Bucket      = "aws_s3_bucket"
	// Source Type
	sourceType         = "source_type"
	fileDomain         = "file_domain"
	RetryBackoff       = "retry_backoff"
	GRPCRequestTimeout = "grpc_request_timeout"
	// Rest
	apiURL                = "api_url"
	restConnectionTimeout = "rest_connection_timeout"
)

//nolint:gochecknoinits // This can be ignored
func init() {
	// Config file setup
	viper.SetConfigName("panorama-fleet-hauler")
	viper.SetConfigType("json")

	// panorama fleet hauler microservice
	viper.SetDefault(httpPort, constants.HTTPPort)
	viper.SetDefault(logLevel, constants.LogLevel)
	viper.SetDefault(ginlogLevel, constants.GinlogLevel)

	viper.SetDefault(config, constants.Local)
	viper.SetDefault(debugMode, false)
	viper.SetDefault(pprofEnabled, false)
	viper.SetDefault(deployMode, constants.DeployModeTest) // Change to Prod

	// grpc
	viper.SetDefault(grpcConnectionTimeout, constants.GRPCConnectionTimeoutSecs)
	// Fleet Grpc Endpoints
	viper.SetDefault(fleetGrpcEndpoint, constants.FleetGrpcEndpoint)
	viper.SetDefault(fleetGrpcDeviceType1Endpoint, constants.FleetGrpcDeviceType1Endpoint)
	viper.SetDefault(fleetGrpcDeviceType2Endpoint, constants.FleetGrpcDeviceType2Endpoint)

	// Authz Endpoint
	viper.SetDefault(authzEndpoint, constants.AuthzEndpoint)
	viper.SetDefault(authzDisable, constants.DisableAuthz)
	viper.SetDefault(updateServiceAuthExpiryTimeoutMins, constants.PFHUpdateAuthzServiceTokenExpiryTimeoutMins)

	// Grip mock
	viper.SetDefault(gripmockConnectionTimeout, constants.GripMockTimeoutSecs)

	// Customer ID
	viper.SetDefault(useCustomerID, constants.UseCustomerID)
	viper.SetDefault(localCluster, constants.LocalCluster)

	// You can change any of these defaults using environment variables or env in deployment yaml
	viper.AutomaticEnv()

	// kafka
	viper.SetDefault(HaulerConsumerKafkaTopic, constants.HaulerConsumerKafkaTopic)
	viper.SetDefault(HaulerProducerKafkaTopic, constants.HaulerProducerKafkaTopic)
	viper.SetDefault(HaulerHarmonyKafkaTopic, constants.HaulerHarmonyKafkaTopic)
	viper.SetDefault(HaulerKafkaGroupID, constants.HaulerConsumerGroupID)

	viper.SetDefault(kafkaBootstrapServers, constants.KafkaBootstrapServers)
	viper.SetDefault(kafkaTimeoutInSeconds, 10) //nolint:gomnd //10 seconds
	viper.SetDefault(kafkaDualStack, true)
	viper.SetDefault(kafkaBatchSize, 1)                // we need a pretty quick response back to CCS, so don't wait for a batch to fill up
	viper.SetDefault(kafkaTLSOn, constants.KafkaTLSOn) // TLS encryption is on for AWS and off for ccs-dev environment
	viper.SetDefault(kafkaEventSource, constants.KafkaEventSource)
	viper.SetDefault(kafkaEventType, constants.KafkaEventType)
	viper.SetDefault(internalJwt, constants.InternalJwt)
	viper.SetDefault(AWSAccessKey, constants.AWSAccessKeyID)
	viper.SetDefault(AWSRegion, constants.AWSRegion)
	viper.SetDefault(AWSScrtAccessKey, constants.AWSSecretAccessKey)
	viper.SetDefault(AWSS3Bucket, constants.AWSS3Bucket)
	viper.SetDefault(sourceType, constants.SourceType)
	viper.SetDefault(fileDomain, constants.FileDomain)
	viper.SetDefault(RetryBackoff, constants.RetryBackoff)
	// Rest
	viper.SetDefault(apiURL, constants.APIURL)
	viper.SetDefault(restConnectionTimeout, constants.RestTimeoutSecs)
}

// GetHTTPPort returns port to listen on for HTTP requests
func GetHTTPPort() string {
	return viper.GetString(httpPort)
}

// GetLogLevel returns the log level
func GetLogLevel() string {
	return viper.GetString(logLevel)
}

func GetGinLogLevel() string {
	return viper.GetString(ginlogLevel)
}

// Config returns the configuration in which the application is running
func GetConfig() string {
	return viper.GetString(config)
}

// DebugMode returns true if debug mode is on
func GetDebugMode() bool {
	return viper.GetBool(debugMode)
}

// DeloyMode returns test/production
func GetDeployMode() string {
	return viper.GetString(deployMode)
}

// GetGRPCConnectionTimeout returns the gRPC server connection timeout
func GetGRPCConnectionTimeout() time.Duration {
	return time.Duration(viper.GetInt(grpcConnectionTimeout)) * time.Second
}

// GetAuthzEndpoint returns the endpoint used to communicate with the SC-Authz service for authorization reqeusts
func GetAuthzEndpoint() string {
	return viper.GetString(authzEndpoint)
}

// GetAuthzEndpoint returns the endpoint used to communicate with the SC-Authz service for authorization reqeusts
func GetAuthzDisableStatus() bool {
	return viper.GetBool(authzDisable)
}

// SetAuthzDisableStatus SetAuthzEndpoint sets the status indicating if SC-Authz service will be used or not for authorization reqeusts
func SetAuthzDisableStatus(authzStatus bool) {
	viper.SetDefault(authzDisable, authzStatus)
}

// GetGripmockConnectionTimeout returns the grip mock server connection timeout
func GetGripMockConnectionTimeout() time.Duration {
	return time.Duration(viper.GetInt(gripmockConnectionTimeout)) * time.Second
}

// GetUseCustomerID returns True or False if Customer ID should be used in GETs call
func GetUseCustomerID() bool {
	return viper.GetBool(useCustomerID)
}

func SetUseCustomerID(flag bool) {
	viper.Set(useCustomerID, flag)
}

func GetFleetGrpcEndpoint() string {
	return viper.GetString(fleetGrpcEndpoint)
}

func SetFleetGrpcEndpoint(fgEndpoint string) {
	viper.Set(fleetGrpcEndpoint, fgEndpoint)
}

func GetFleetGrpcDeviceType1Endpoint() string {
	return viper.GetString(fleetGrpcDeviceType1Endpoint)
}

func SetFleetGrpcDeviceType1Endpoint(fgEndpoint string) {
	viper.Set(fleetGrpcDeviceType1Endpoint, fgEndpoint)
}

func GetFleetGrpcDeviceType2Endpoint() string {
	return viper.GetString(fleetGrpcDeviceType2Endpoint)
}

func SetFleetGrpcDeviceType2Endpoint(fgEndpoint string) {
	viper.Set(fleetGrpcDeviceType2Endpoint, fgEndpoint)
}

// Get GetKafkaEventSource returns kafka event source
func GetKafkaEventSource() string {
	return viper.GetString(kafkaEventSource)
}

func GetKafkaEventType() string {
	return viper.GetString(kafkaEventType)
}

// GetHaulerKafkaGroupID returns consumer group id
func GetHaulerKafkaGroupID() string {
	return viper.GetString(HaulerKafkaGroupID)
}

// GetHaulerConsumerKafkaTopic returns the topic to use for Fleet Hauler
func GetHaulerConsumerKafkaTopic() string {
	return viper.GetString(HaulerConsumerKafkaTopic)
}

// GetHaulerProducerKafkaTopic returns the topic to use for Fleet Hauler
func GetHaulerProducerKafkaTopic() string {
	return viper.GetString(HaulerProducerKafkaTopic)
}

// GetHaulerHarmonyKafkaTopic returns the topic to use for Fleet Hauler
func GetHaulerHarmonyKafkaTopic() string {
	return viper.GetString(HaulerHarmonyKafkaTopic)
}

// GetKafkaBootstrapServers returns a comma separated list of Kafka bootstrap servers
func GetKafkaBootstrapServers() string {
	return viper.GetString(kafkaBootstrapServers)
}

// How long to wait for a connection to the Kafka broker before giving up
func GetKafkaTimeoutInSeconds() int {
	return viper.GetInt(kafkaTimeoutInSeconds)
}

// Should we use dual IPv4 and IPv6 stack?
func GetKafkaDualStackOn() bool {
	return viper.GetBool(kafkaDualStack)
}

func GetKafkaBatchSize() int {
	return viper.GetInt(kafkaBatchSize)
}

// Should we use secure mode when connecting to kafka?
func GetKafkaTLSOn() bool {
	return viper.GetBool(kafkaTLSOn)
}

func GetInternalJwt() string {
	return viper.GetString(internalJwt)
}

func GetUpdateAuthzServiceTokenTimeout() uint32 {
	return viper.GetUint32(updateServiceAuthExpiryTimeoutMins)
}

// GetAWSRegions returns AWS account regions
func GetAWSRegion() string {
	return viper.GetString(AWSRegion)
}

// GetAWSRegions returns AWS account regions
func GetAWSAccessKey() string {
	return viper.GetString(AWSAccessKey)
}

// GetAWSRegions returns AWS account regions
func GetAWSSecretAccessKey() string {
	return viper.GetString(AWSScrtAccessKey)
}

// GetAWSS3BucketName returns AWS bucket name
func GetAWSS3BucketName() string {
	return viper.GetString(AWSS3Bucket)
}

func GetPprofEnabled() bool {
	return viper.GetBool(pprofEnabled)
}

// GetSourceType returns Source type
func GetSourceType() string {
	return viper.GetString(sourceType)
}

func GetFileDomain() string {
	return viper.GetString(fileDomain)
}

func GetRetryBackoff() string {
	return viper.GetString(RetryBackoff)
}

func GetRestConnectionTimeout() time.Duration {
	return time.Duration(viper.GetInt(restConnectionTimeout)) * time.Second
}

func GetLocalCluster() bool {
	return viper.GetBool(localCluster)
}

func SetLocalCluster(flag bool) {
	viper.Set(localCluster, flag)
}

func GetAPIURL() string {
	return viper.GetString(apiURL)
}
