// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

package restclient_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/restclient"
)

const (
	mockToken = "token"
)

type RequestHandler func(w http.ResponseWriter, r *http.Request)

type RestClientSuite struct {
	suite.Suite
	mockRestServer    *httptest.Server
	mockRestServerURL *url.URL

	mockRestSrv    *httptest.Server
	mockRestSrvURL *url.URL
}

func (rcSuite *RestClientSuite) SetupSuite() {
	rcSuite.mockRestServer = createServerMock(createDummyHandlerMock)

	rcSuite.mockRestServerURL, _ = url.Parse(rcSuite.mockRestServer.URL)

	rcSuite.mockRestSrv = createServerMock(timeoutMock)

	rcSuite.mockRestSrvURL, _ = url.Parse(rcSuite.mockRestSrv.URL)
}

func (rcSuite *RestClientSuite) TearDownSuite() {
	rcSuite.mockRestServer.Close()
}

func (rcSuite *RestClientSuite) TestClientDummyCreate() {
	headers := make(map[string]string)
	headers["X-Auth-Token"] = mockToken

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/dummy", bytes.NewReader([]byte(`{"name":"mockDummy"}`)))
	rcSuite.NoError(err)

	client, _ := restclient.NewRestClient("http://localhost:"+rcSuite.mockRestServerURL.Port(), 5*time.Second)
	req.Header.Set("X-Auth-Token", mockToken)
	resp, err := client.SendRequest(context.TODO(), req)
	rcSuite.NoError(err)
	rcSuite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	rcSuite.NoError(err)
	rcSuite.Equal("mock dummy created", string(bodyBytes))
}

func (rcSuite *RestClientSuite) TestClientBaseUrl() {
	headers := make(map[string]string)
	headers["X-Auth-Token"] = mockToken

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/dummy", bytes.NewReader([]byte(`{"name":"mockDummy"}`)))
	rcSuite.NoError(err)

	client, _ := restclient.NewRestClient("", 5*time.Second)
	req.Header.Set("X-Auth-Token", mockToken)
	client.SetBaseURL("http", "localhost", rcSuite.mockRestServerURL.Port())

	resp, err := client.SendRequest(context.TODO(), req)
	rcSuite.NoError(err)
	rcSuite.Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	rcSuite.NoError(err)
	rcSuite.Equal("mock dummy created", string(bodyBytes))
}

func (rcSuite *RestClientSuite) TestClientDoErr() {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/dummy", bytes.NewReader([]byte(`{"name":"mockDummy"}`)))
	rcSuite.NoError(err)

	// try connect to 5392 which no-one is listening in DVM
	client, _ := restclient.NewRestClient("http://localhost:5392", 5*time.Second)
	resp, connectErr := client.SendRequest(context.TODO(), req)
	if connectErr != nil {
		rcSuite.Error(connectErr)
	} else {
		defer resp.Body.Close()
	}
}

func (rcSuite *RestClientSuite) TestClientTimeout() {
	headers := make(map[string]string)
	headers["X-Auth-Token"] = mockToken

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/dummy", bytes.NewReader([]byte(`{"name":"mockDummy"}`)))
	rcSuite.NoError(err)

	client, _ := restclient.NewRestClient("http://localhost:"+rcSuite.mockRestSrvURL.Port(), 5*time.Second)
	req.Header.Set("X-Auth-Token", mockToken)
	resp, ctxErr := client.SendRequest(context.TODO(), req)
	if ctxErr != nil {
		rcSuite.Error(ctxErr)
		rcSuite.Contains(ctxErr.Error(), "context deadline exceeded")
	} else {
		defer resp.Body.Close()
	}
}

func TestRestClientSuite(t *testing.T) {
	suite.Run(t, new(RestClientSuite))
}

func createServerMock(reqHandler RequestHandler) *httptest.Server {
	handler := http.HandlerFunc(reqHandler)
	srv := httptest.NewServer(handler)
	return srv
}

func createDummyHandlerMock(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("mock dummy created"))
}

func timeoutMock(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
}
