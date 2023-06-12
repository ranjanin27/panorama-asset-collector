// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package configs

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const jsonStream = `{
	"kafka_bootstrap_servers": "0.0.0.0",
	"gripmock_connection_timeout": "30"
	}`
const (
	deviceType1EndPoint = "http://127.0.0.1:5565"
	deviceType2EndPoint = "http://127.0.0.1:5565"
	fleetEndPoint       = "http://127.0.0.1:5565"
)

func TestConfiguration(t *testing.T) {
	r := strings.NewReader(jsonStream)
	viper.SetConfigType("json")
	err := viper.ReadConfig(r)
	assert.Nil(t, err, "Could not read test config")

	assert.False(t, GetUseCustomerID())
	SetUseCustomerID(true)
	assert.True(t, GetUseCustomerID())
	assert.NotZero(t, GetConfig())
	assert.NotEmpty(t, GetHTTPPort())
	assert.NotZero(t, GetLogLevel())
	assert.NotZero(t, GetGinLogLevel())
	assert.False(t, GetDebugMode())
	assert.NotZero(t, GetDeployMode())
	assert.NotEmpty(t, GetGRPCConnectionTimeout())

	assert.NotZero(t, GetAuthzEndpoint())
	assert.NotEmpty(t, GetGripMockConnectionTimeout())

	assert.NotEmpty(t, GetHaulerConsumerKafkaTopic())
	assert.NotEmpty(t, GetHaulerProducerKafkaTopic())
	assert.NotEmpty(t, GetHaulerHarmonyKafkaTopic())
	assert.NotEmpty(t, GetKafkaBootstrapServers())
	assert.NotEmpty(t, GetKafkaEventSource())
	assert.NotEmpty(t, GetKafkaEventType())
	assert.Empty(t, GetInternalJwt())
	assert.NotZero(t, GetUpdateAuthzServiceTokenTimeout())
	GetKafkaTLSOn()
	GetKafkaTimeoutInSeconds()
	GetKafkaBatchSize()
	GetKafkaDualStackOn()
	GetHaulerKafkaGroupID()
	GetAWSAccessKey()
	GetAWSRegion()
	GetAWSSecretAccessKey()
	GetAWSS3BucketName()
	GetPprofEnabled()
	GetSourceType()
	GetFileDomain()
	GetRetryBackoff()
	GetAPIURL()
	SetFleetGrpcEndpoint(fleetEndPoint)
	GetRestConnectionTimeout()
	assert.Equal(t, GetFleetGrpcEndpoint(), fleetEndPoint)
	SetFleetGrpcDeviceType1Endpoint(deviceType1EndPoint)
	assert.Equal(t, GetFleetGrpcDeviceType1Endpoint(), deviceType1EndPoint)
	SetFleetGrpcDeviceType2Endpoint(deviceType2EndPoint)
	assert.Equal(t, GetFleetGrpcDeviceType2Endpoint(), deviceType2EndPoint)
	SetAuthzDisableStatus(true)
	assert.Equal(t, GetAuthzDisableStatus(), true)
	SetLocalCluster(true)
	assert.Equal(t, GetLocalCluster(), true)
}
