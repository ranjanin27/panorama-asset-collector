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
	"github.hpe.com/nimble-dcs/panorama-s3-upload/pkg/s3"
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

type AssetInterface interface {
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
	GetProtectionStores(ctx context.Context, authHeader string) ([]model.ProtectionStore, error)
	GetProtectionStoreGateways(ctx context.Context, authHeader string) ([]model.ProtectionStoreGateway, error)
	GetStoreonces(ctx context.Context, authHeader string) ([]model.Storeonce, error)
	GetCSPAccounts(ctx context.Context, authHeader string) ([]model.CSPAccount, error)
	GetCSPVolumes(ctx context.Context, authHeader string) ([]model.CSPVolume, error)
	GetDOs(ctx context.Context, authHeader string) ([]model.DO, error)
	GetMsSqlDB(ctx context.Context, authHeader string) ([]model.MsSqlDB, error)
	GetMsSqlInstances(ctx context.Context, authHeader string) ([]model.MsSqlInstance, error)
	GetDBBackups(ctx context.Context, dbId, authHeader string) ([]model.MsSqlDBBackup, error)
	GetDBSnapshots(ctx context.Context, dbId, authHeader string) ([]model.MsSqlDBSnapshot, error)
	GetMsSqlProtectionGroups(ctx context.Context, authHeader string) ([]model.MsSqlProtectionGroup, error)
}

type CommonClient struct {
	client     restclient.RestInterface
	ctx        context.Context
	customerID string
}

func NewCommonClient(ctx context.Context) (AssetInterface, error) {
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

func (assetClient *CommonClient) HandleRequest(ctx context.Context, req *http.Request,
	headers map[string]string) (response *http.Response, status int, err error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, sendRequestErr := assetClient.client.SendRequest(ctx, req)
	if sendRequestErr != nil {
		logger.WithContext(ctx).Error(sendRequestErr.Error())
		return nil, -1, sendRequestErr
	}
	// defer resp.Body.Close() - DONT close here, close at callbacks once we are done with building msg
	return resp, resp.StatusCode, nil
}

func (assetClient *CommonClient) SetCustomerIDForRest(customerID string) {
	assetClient.customerID = customerID
}

func (assetClient *CommonClient) GetAuthHeaderForRest() (string, error) {
	urlStr := "https://sso.common.cloud.hpe.com/as/token.oauth2"

	// HCIPOC account
	//payload := strings.NewReader("grant_type=client_credentials&client_id=0d6f9b5d-c528-4826-9340-f21a6500d960&client_secret=9a2b76567aef11edb5397eb97d380c5e")
	// SCDEV01 account
	payload := strings.NewReader("grant_type=client_credentials&client_id=7fb4cf66-af63-4953-94c4-12d9e63b081c&client_secret=2d6939bc142e11ee9ad40667a7649451")

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

var UploadToS3 = func(ctx context.Context, jsonContent *[]byte, bucketName, awsS3Region, awsAccessKeyID,
	awsSecretAccessKey string, uploadType s3.UploadType) (string, int, error) {
	logger.Debugf("bucketName: %v awsS3Region %v awsAccessKeyID %v awsSecretAccessKey %v uploadType %v\n",
		bucketName, awsS3Region, awsAccessKeyID, awsSecretAccessKey, s3.UploadType(configs.GetSourceType()))
	keyname, filesize, err := s3.UploadToAwsS3(jsonContent, bucketName, awsS3Region, awsAccessKeyID,
		awsSecretAccessKey, s3.UploadType(configs.GetSourceType()))
	return keyname, filesize, err
}

func (assetClient *CommonClient) GetVMs(ctx context.Context, authHeader string) ([]model.VirtualMachine,
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		vmData, err := json.Marshal(vms)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &vmData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))*/
	// allItems now contains all the astorageSystems from the paginated API
	return vms, nil
}

func (assetClient *CommonClient) GetDatastores(ctx context.Context, authHeader string) ([]model.Datastore,
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		dsData, err := json.Marshal(datastores)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &dsData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return datastores, nil
}

func (assetClient *CommonClient) GetProtectionPolicies(ctx context.Context, authHeader string) ([]model.ProtectionPolicy,
	error) {
	baseURL := "/backup-recovery/v1beta1/protection-policies"
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		ppData, err := json.Marshal(protectionpolicies)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &ppData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return protectionpolicies, nil
}

func (assetClient *CommonClient) GetVMProtectionGroups(ctx context.Context, authHeader string) ([]model.VMProtectionGroup,
	error) {
	baseURL := "/backup-recovery/v1beta1/virtual-machine-protection-groups"
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		vmpgData, err := json.Marshal(vmpg)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &vmpgData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return vmpg, nil
}

func (assetClient *CommonClient) GetVMBackups(ctx context.Context, vmId, authHeader string) ([]model.VMBackup, error) {
	baseURL := "/backup-recovery/v1beta1/virtual-machines/" + vmId + "/backups"
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
		for i := 0; i < len(paginatedResponse.Items); i++ {
			paginatedResponse.Items[i].SourceID = vmId
		}
		vmbkps = append(vmbkps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(vmbkps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	/*
		vmbkpsData, err := json.Marshal(vmbkps)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &vmbkpsData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return vmbkps, nil
}

func (assetClient *CommonClient) GetVMSnapshots(ctx context.Context, vmId, authHeader string) ([]model.VMSnapshot, error) {
	baseURL := "/backup-recovery/v1beta1/virtual-machines/" + vmId + "/snapshots"
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		vmData, err := json.Marshal(vmsnaps)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &vmData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return vmsnaps, nil
}

func (assetClient *CommonClient) GetDSBackups(ctx context.Context, dsId, authHeader string) ([]model.DatastoreBackup, error) {
	baseURL := "/backup-recovery/v1beta1/datastores/" + dsId + "/backups"
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
		for i := 0; i < len(paginatedResponse.Items); i++ {
			paginatedResponse.Items[i].SourceID = dsId
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
	/*
		dsbkpsData, err := json.Marshal(dsbkps)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &dsbkpsData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return dsbkps, nil
}

func (assetClient *CommonClient) GetDSSnapshots(ctx context.Context, dsId, authHeader string) ([]model.DSSnapshot, error) {
	baseURL := "/backup-recovery/v1beta1/datastores/" + dsId + "/snapshots"
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		dsbkpsData, err := json.Marshal(dssnaps)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &dsbkpsData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return dssnaps, nil
}

func (assetClient *CommonClient) GetProtectedVMs(ctx context.Context, authHeader string) ([]model.ProtectedVM, error) {
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		pvmsData, err := json.Marshal(pvms)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &pvmsData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return pvms, nil
}

func (assetClient *CommonClient) GetCSPMachineInstances(ctx context.Context, authHeader string) ([]model.CSPMachineInstance, error) {
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		miData, err := json.Marshal(cspmis)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &miData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return cspmis, nil
}

func (assetClient *CommonClient) GetZertoVPGs(ctx context.Context, authHeader string) ([]model.ZertoVPG, error) {
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
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
	/*
		zvData, err := json.Marshal(zvpgs)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &zvData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return zvpgs, nil
}

func (assetClient *CommonClient) GetProtectionStores(ctx context.Context, authHeader string) ([]model.ProtectionStore, error) {
	baseURL := "/backup-recovery/v1beta1/protection-stores"
	var ps []model.ProtectionStore
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.ProtectionStores
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		ps = append(ps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(ps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	/*
		psData, err := json.Marshal(ps)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &psData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return ps, nil
}

func (assetClient *CommonClient) GetProtectionStoreGateways(ctx context.Context, authHeader string) ([]model.ProtectionStoreGateway, error) {
	baseURL := "/api/v1/protection-store-gateways"
	//baseURL := "/backup-recovery/v1beta1/protection-store-gateways"
	var psgs []model.ProtectionStoreGateway
	// create the request with the appropriate URL and query params
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, bytes.NewReader(nil))
	if reqErr != nil {
		logger.WithContext(ctx).Error(reqErr.Error())
		return nil, reqErr
	}
	req.Close = true
	// add authorization header to the req
	if !configs.GetLocalCluster() {
		req.Header.Add(RestAuthHeader, fmt.Sprintf("Bearer %s", authHeader))
	}
	// make the HTTP request
	resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
	if handleReqErr != nil {
		logger.WithContext(ctx).Error(handleReqErr.Error())
		return nil, handleReqErr
	}
	// decode the JSON response into a PaginatedResponse object
	var paginatedResponse model.ProtectionStoreGateways
	decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
	if decodeErr != nil {
		logger.WithContext(ctx).Error(decodeErr.Error())
		resp.Body.Close()
		return nil, decodeErr
	}

	// append the items from the current page to the astoragePools slice
	psgs = append(psgs, paginatedResponse.Items...)
	/*
		psData, err := json.Marshal(psgs)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &psData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return psgs, nil
}

func (assetClient *CommonClient) GetStoreonces(ctx context.Context, authHeader string) ([]model.Storeonce, error) {
	baseURL := "/backup-recovery/v1beta1/storeonces"
	var sos []model.Storeonce
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.Storeonces
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		sos = append(sos, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(sos) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	/*
		psData, err := json.Marshal(sos)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &psData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
		// allItems now contains all the astorageSystems from the paginated API
	*/
	return sos, nil
}

func (assetClient *CommonClient) GetCSPAccounts(ctx context.Context, authHeader string) ([]model.CSPAccount, error) {
	baseURL := "/api/v1/csp-accounts"
	var cspa []model.CSPAccount
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.CSPAccounts
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		cspa = append(cspa, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(cspa) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	/*
		cspData, err := json.Marshal(cspa)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &cspData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
		// allItems now contains all the astorageSystems from the paginated API
	*/
	return cspa, nil
}

func (assetClient *CommonClient) GetCSPVolumes(ctx context.Context, authHeader string) ([]model.CSPVolume, error) {
	baseURL := "/api/v1/csp-volumes"
	var cspv []model.CSPVolume
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.CSPVolumes
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		cspv = append(cspv, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(cspv) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}

	cspData, err := json.Marshal(cspv)
	if err != nil {
		logger.WithContext(ctx).Error("marshalling error %s", err)
		return nil, err
	}
	UploadToS3(ctx, &cspData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
		configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))

	// allItems now contains all the astorageSystems from the paginated API
	return cspv, nil
}

func (assetClient *CommonClient) GetDOs(ctx context.Context, authHeader string) ([]model.DO, error) {
	baseURL := "/backup-recovery/v1beta1/data-orchestrators"
	var do []model.DO
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.DOs
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		do = append(do, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(do) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	/*
		cspData, err := json.Marshal(cspa)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &cspData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
		// allItems now contains all the astorageSystems from the paginated API
	*/
	return do, nil
}

func (assetClient *CommonClient) GetMsSqlDB(ctx context.Context, authHeader string) ([]model.MsSqlDB, error) {
	baseURL := "/backup-recovery/v1beta1/mssql-databases"
	var db []model.MsSqlDB
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.MsSqlDBs
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		db = append(db, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(db) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	return db, nil
}
func (assetClient *CommonClient) GetMsSqlInstances(ctx context.Context, authHeader string) ([]model.MsSqlInstance, error) {
	baseURL := "/backup-recovery/v1beta1/mssql-instances"
	var dbIns []model.MsSqlInstance
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.MsSqlInstances
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		dbIns = append(dbIns, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(dbIns) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	return dbIns, nil
}

func (assetClient *CommonClient) GetDBBackups(ctx context.Context, dbId, authHeader string) ([]model.MsSqlDBBackup, error) {
	baseURL := "/backup-recovery/v1beta1/mssql-databases/" + dbId + "/backups"
	var dbbkps []model.MsSqlDBBackup
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.MsSqlDBBackups
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		for i := 0; i < len(paginatedResponse.Items); i++ {
			paginatedResponse.Items[i].SourceID = dbId
		}
		dbbkps = append(dbbkps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(dbbkps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	/*
		vmbkpsData, err := json.Marshal(vmbkps)
		if err != nil {
			logger.WithContext(ctx).Error("marshalling error %s", err)
			return nil, err
		}
		UploadToS3(ctx, &vmbkpsData, configs.GetAWSS3BucketName(), configs.GetAWSRegion(),
	                configs.GetAWSAccessKey(), configs.GetAWSSecretAccessKey(), s3.UploadType(configs.GetSourceType()))
	*/
	// allItems now contains all the astorageSystems from the paginated API
	return dbbkps, nil
}

func (assetClient *CommonClient) GetDBSnapshots(ctx context.Context, dbId, authHeader string) ([]model.MsSqlDBSnapshot, error) {
	baseURL := "/backup-recovery/v1beta1/mssql-databases/" + dbId + "/snapshots"
	var dbsnaps []model.MsSqlDBSnapshot
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.MsSqlDBSnapshots
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		for i := 0; i < len(paginatedResponse.Items); i++ {
			paginatedResponse.Items[i].SourceID = dbId
		}
		dbsnaps = append(dbsnaps, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(dbsnaps) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}
	return dbsnaps, nil
}
func (assetClient *CommonClient) GetMsSqlProtectionGroups(ctx context.Context, authHeader string) ([]model.MsSqlProtectionGroup, error) {
	baseURL := "/backup-recovery/v1beta1/mssql-database-protection-groups"
	var dbPG []model.MsSqlProtectionGroup
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
		resp, _, handleReqErr := assetClient.HandleRequest(ctx, req, nil)
		if handleReqErr != nil {
			logger.WithContext(ctx).Error(handleReqErr.Error())
			return nil, handleReqErr
		}
		// decode the JSON response into a PaginatedResponse object
		var paginatedResponse model.MsSqlProtectionGroups
		decodeErr := json.NewDecoder(resp.Body).Decode(&paginatedResponse)
		if decodeErr != nil {
			logger.WithContext(ctx).Error(decodeErr.Error())
			resp.Body.Close()
			return nil, decodeErr
		}
		// append the items from the current page to the astoragePools slice
		dbPG = append(dbPG, paginatedResponse.Items...)
		// break the loop if we have fetched all the astoragePools
		if len(dbPG) >= paginatedResponse.Total {
			resp.Body.Close()
			break
		}
		// update the page offset for the next iteration
		pageOffset += pageLimit
	}

	return dbPG, nil
}
