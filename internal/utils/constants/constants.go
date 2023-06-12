// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package constants

// const are exported. Configs for kafka and db
const (
	// Rest protocol
	HTTPS = "https"
	HTTP  = "http"

	// PORTs
	HTTPPort = "8300"

	Local          = "local"
	Localhost      = "localhost"
	DeployModeTest = "test"
	DeployModeProd = "prod"

	LogLevel         = "debug"
	GinlogLevel      = "release"
	TraceStatePrefix = "b3=0"

	// Authz Endpoint
	AuthzEndpoint                               = "grpc.authz.svc.cluster.local:5100"
	PFHUpdateAuthzServiceTokenExpiryTimeoutMins = 120

	// Authz Disable Status
	DisableAuthz = false

	// grip mock configs
	GripMockTimeoutSecs = 30

	// Customer ID
	UseCustomerID = false
	LocalCluster  = false
	// auth for jwt
	AuthorizationHeader = "Authorization"

	// GRPC
	GRPCConnectionTimeoutSecs = 10
	RetryBackoff              = "3s"

	// Fleet Grpc Endpoint
	FleetGrpcEndpoint            = "fleet-grpc-api-service.fleet.svc.cluster.local:5565"
	FleetGrpcDeviceType1Endpoint = "fleet-grpc-api-device-type1-service.fleet.svc.cluster.local:5565"
	FleetGrpcDeviceType2Endpoint = "fleet-grpc-api-device-type2-service.fleet.svc.cluster.local:5565"
	FleetNbRestURL               = "fleet-nb-rest.fleet.svc.cluster.local:3443"

	// kafka configs
	KafkaBootstrapServers = "ccs-kafka:19092"
	KafkaTLSOn            = false

	// Actual Topic Name
	HaulerConsumerKafkaTopic = "panorama.common.hauler.collection"
	HaulerProducerKafkaTopic = "panorama.hauler.collection.status"
	HaulerHarmonyKafkaTopic  = "panorama.common.hdpfilesync"
	HaulerConsumerGroupID    = "panorama.common.hauler.consumergroup"

	KafkaEventSource = "https://github.hpe.com/nimble-dcs/panorama-common-hauler"
	KafkaEventType   = "fleet.common.hauler.collection"

	// JSON Version
	Version = "1.0"

	// Internal JWT token
	InternalJwt = ""

	// prometheus metrics
	PrometheusMetrics = "panorama_common_hauler"

	// AWS related
	AWSRegion          = ""
	AWSAccessKeyID     = ""
	AWSSecretAccessKey = ""
	AWSS3Bucket        = ""

	// Source type
	SourceType = "Fleet"
	FileDomain = "panorama"

	// Rest
	RestTimeoutSecs = 60
	APIURL          = "https://hcipoc-app.qa.cds.hpe.com"
)
