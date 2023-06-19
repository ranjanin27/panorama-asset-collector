// (C) Copyright 2023 Hewlett Packard Enterprise Development LP

//nolint:dupl // Ignore
package commonclient

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/commonclient/model"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/clients/restclient"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/configs"
	"github.hpe.com/nimble-dcs/panorama-common-hauler/internal/utils/logging"
)

const (
	RestAuthHeader               = "Authorization"
	PFHAuthzImplicitPermissionID = "a8af5fd6-8935-40f3-9e38-5438eebfe771"
	DeviceType4Tag               = "device-type4"
	numAuthHeaders               = 5
	TokenDir                     = "/tmp"
	PageLimit                    = 50
	PageOffset                   = 0
	baseURI                      = "/api/v1/storage-systems"
)

var commonClient *CommonClient

var (
	logger = logging.GetLogger()
	json   = jsoniter.ConfigCompatibleWithStandardLibrary
)

type ArcusInterface interface {
	HandleRequest(context.Context, *http.Request, map[string]string) (*http.Response, int, error)
	GetAuthHeaderForRest() (string, error)
	SetCustomerIDForRest(string)

	GetVMs(ctx context.Context, authHeader string) ([]model.VirtualMachine, error)
	GetDatastores(ctx context.Context, authHeader string) ([]model.Datastore, error)
	GetProtectionPolicies(ctx context.Context, authHeader string) ([]model.ProtectionPolicy, error)
	GetVMProtectionGroups(ctx context.Context, authHeader string) ([]model.VMProtectionGroup, error)
	GetVMBackups(ctx context.Context, vmId, authHeader string) ([]model.VMBackup, error)
	GetVMSnapshots(ctx context.Context, vmId, authHeader string) ([]model.VMSnapshot, error)
	GetDSBackups(ctx context.Context, dsId, authHeader string) ([]model.DatastoreBackup, error)
	GetDSSnapshots(ctx context.Context, dsId, authHeader string) ([]model.DSSnapshot, error)
	GetProtectedVMs(ctx context.Context, authHeader string) ([]model.ProtectedVM, error)
	GetCSPMachineInstances(ctx context.Context, authHeader string) ([]model.CSPMachineInstance, error)
	GetZertoVPGs(ctx context.Context, authHeader string) ([]model.ZertoVPG, error)
}

type CommonClient struct {
	client     restclient.RestInterface
	ctx        context.Context
	customerID string
}

func NewCommonClient(ctx context.Context) (ArcusInterface, error) {
	if commonClient == nil {
		client, err := restclient.NewRestClient(configs.GetAPIURL(), configs.GetRestConnectionTimeout())
		if err != nil {
			logger.WithContext(ctx).Errorf("Restclient creation failed %v", err.Error())
			return nil, err
		}
		logger.WithContext(ctx).Info("Created Arcus Client instance")
		commonClient = &CommonClient{
			client: client,
			ctx:    ctx,
		}
	}
	return commonClient, nil
}

func (arcusClient *CommonClient) HandleRequest(ctx context.Context, req *http.Request,
	headers map[string]string) (response *http.Response, status int, err error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, sendRequestErr := arcusClient.client.SendRequest(ctx, req)
	if sendRequestErr != nil {
		logger.WithContext(ctx).Error(sendRequestErr.Error())
		return nil, -1, sendRequestErr
	}
	// defer resp.Body.Close() - DONT close here, close at callbacks once we are done with building msg
	return resp, resp.StatusCode, nil
}

func (arcusClient *CommonClient) SetCustomerIDForRest(customerID string) {
	arcusClient.customerID = customerID
}

func (arcusClient *CommonClient) GetAuthHeaderForRest() (string, error) {
	urlStr := "https://sso.common.cloud.hpe.com/as/token.oauth2"

	payload := strings.NewReader("grant_type=client_credentials&client_id=0d6f9b5d-c528-4826-9340-f21a6500d960&client_secret=9a2b76567aef11edb5397eb97d380c5e")

	req, err := http.NewRequest("POST", urlStr, payload)
	if err != nil {
		logger.Errorf("Failed to create HTTP request: %v", err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Close = true
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("Failed to send HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string  `json:"access_token"`
		TokenType   string  `json:"token_type"`
		ExpiresIn   float64 `json:"expires_in"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		logger.Errorf("failed to decode token response: %v", err)
		return "", err
	}
	return tokenResponse.AccessToken, nil
}

func (arcusClient *CommonClient) GetVMs(ctx context.Context, authHeader string) ([]model.VirtualMachine,
	error) {
	baseURL := "/api/v1/virtual-machines"
	var vms []model.VirtualMachine
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.VirtualMachines
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		vms = append(vms, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(vms) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return vms, nil
}

func (arcusClient *CommonClient) GetDatastores(ctx context.Context, authHeader string) ([]model.Datastore,
	error) {
	baseURL := "/api/v1/datastores"
	var datastores []model.Datastore
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.Datastores
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		datastores = append(datastores, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(datastores) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return datastores, nil
}

func (arcusClient *CommonClient) GetProtectionPolicies(ctx context.Context, authHeader string) ([]model.ProtectionPolicy,
	error) {
	baseURL := "/api/v1/protection-policies"
	var protectionpolicies []model.ProtectionPolicy
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.ProtectionPolicies
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		protectionpolicies = append(protectionpolicies, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(protectionpolicies) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return protectionpolicies, nil
}

func (arcusClient *CommonClient) GetVMProtectionGroups(ctx context.Context, authHeader string) ([]model.VMProtectionGroup,
	error) {
	baseURL := "/api/v1/virtual-machine-protection-groups"
	var vmpg []model.VMProtectionGroup
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.VMProtectionGroups
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		vmpg = append(vmpg, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(vmpg) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return vmpg, nil
}

func (arcusClient *CommonClient) GetVMBackups(ctx context.Context, vmId, authHeader string) ([]model.VMBackup, error) {
	baseURL := "/api/v1/virtual-machines/" + vmId + "/backups"
	var vmbkps []model.VMBackup
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.VMBackups
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		vmbkps = append(vmbkps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(vmbkps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return vmbkps, nil
}

func (arcusClient *CommonClient) GetVMSnapshots(ctx context.Context, vmId, authHeader string) ([]model.VMSnapshot, error) {
	baseURL := "/api/v1/virtual-machines/" + vmId + "/snapshots"
	var vmsnaps []model.VMSnapshot
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.VMSnapshots
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		vmsnaps = append(vmsnaps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(vmsnaps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return vmsnaps, nil
}

func (arcusClient *CommonClient) GetDSBackups(ctx context.Context, dsId, authHeader string) ([]model.DatastoreBackup, error) {
	baseURL := "/api/v1/datastores/" + dsId + "/backups"
	var dsbkps []model.DatastoreBackup
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.DatastoreBackups
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		dsbkps = append(dsbkps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(dsbkps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return dsbkps, nil
}

func (arcusClient *CommonClient) GetDSSnapshots(ctx context.Context, dsId, authHeader string) ([]model.DSSnapshot, error) {
	baseURL := "/api/v1/datastores/" + dsId + "/snapshots"
	var dssnaps []model.DSSnapshot
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.DSSnapshots
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		dssnaps = append(dssnaps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(dssnaps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return dssnaps, nil
}

func (arcusClient *CommonClient) GetProtectedVMs(ctx context.Context, authHeader string) ([]model.ProtectedVM, error) {
	baseURL := "/disaster-recovery/v1beta1/protected-vms"
	var pvms []model.ProtectedVM
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.ProtectedVMs
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		pvms = append(pvms, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(pvms) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return pvms, nil
}

func (arcusClient *CommonClient) GetCSPMachineInstances(ctx context.Context, authHeader string) ([]model.CSPMachineInstance, error) {
	baseURL := "/api/v1/csp-machine-instances"
	var cspmis []model.CSPMachineInstance
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.CSPMachineInstances
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		cspmis = append(cspmis, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(cspmis) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return cspmis, nil
}

func (arcusClient *CommonClient) GetZertoVPGs(ctx context.Context, authHeader string) ([]model.ZertoVPG, error) {
	baseURL := "/disaster-recovery/v1beta1/virtual-continuous-protection-groups"
	var zvpgs []model.ZertoVPG
	pageLimit := PageLimit   // set the page limit
	pageOffset := PageOffset // set the initial page offset

	for {
		// create the request with the appropriate URL and query params
		req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
		if reqErr != nil {
			logger.WithContext(ctx).Error(reqErr.Error())
			return nil, reqErr
		}
		q := req.URL.Query()
		q.Add("limit", fmt.Sprintf("%d", pageLimit))
		q.Add("offset", fmt.Sprintf("%d", pageOffset))
		req.URL.RawQuery = q.Encode()
		req.Close = true
		// add authorization header to the req
		if !configs.GetLocalCluster() {
			req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
		}
		// make the HTTP request
		resp, _, handleReqErr := arcusClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.ZertoVPGs
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		zvpgs = append(zvpgs, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(zvpgs) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	// allItems now contains all the astorageSystems from the paginated API
	return zvpgs, nil
}
