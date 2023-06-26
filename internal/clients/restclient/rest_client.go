// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

package restclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
)

var (
	logger = logging.GetLogger()
)

type RestInterface interface {
	SetBaseURL(string, string, string)
	SendRequest(ctx context.Context, request *http.Request) (*http.Response, error)
}

type RestClient struct {
	client  *http.Client
	baseURL string
}

func NewRestClient(hostName string, timeout time.Duration) (RestInterface, error) {
	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse("http://web-proxy.corp.hpecorp.net:8080")
		},
		TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
	}
	client := &http.Client{Transport: tr, Timeout: timeout}

	restClient := &RestClient{
		client:  client,
		baseURL: hostName,
	}
	return restClient, nil
}

func (client *RestClient) SetBaseURL(protocol, hostName, port string) {
	client.baseURL = fmt.Sprintf("%s://%s:%s", protocol, hostName, port)
}

func (client *RestClient) SendRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	path := fmt.Sprintf("%s%s", client.baseURL, request.URL.String())
	u, err := url.Parse(path)
	if err != nil {
		logger.WithContext(ctx).Error("Error parsing url : ", err.Error())
		return nil, err
	}
	request.URL = u

	logger.WithContext(ctx).Debug("Request : ", request.Method, " ", request.URL.String())
	request.Header.Set("Content-Type", "application/json")

	response, err := client.client.Do(request)
	if err != nil {
		logger.WithContext(ctx).Errorf("error in request : %v", err.Error())
		return nil, err
	}

	// defer response.Body.Close() - DONT close here, close at callbacks once we are done with building msg
	logger.WithContext(ctx).Debugf("Request processed successfully with response code : %v", response.Status)
	return response, nil
}
