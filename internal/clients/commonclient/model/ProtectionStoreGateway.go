package model

type ProtectionStoreGateway struct {
	CustomerID   string `json:"customerId"`
	Generation   int    `json:"generation"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	ResourceURI  string `json:"resourceUri"`
	ConsoleURI   string `json:"consoleUri"`
	Type         string `json:"type"`
	DatastoreIds []struct {
		DatastoreID string `json:"datastoreId"`
	} `json:"datastoreIds"`
	DatastoresInfo []struct {
		ID                      string `json:"id"`
		Type                    string `json:"type"`
		ResourceURI             string `json:"resourceUri"`
		TotalProvisionedDiskTiB int    `json:"totalProvisionedDiskTiB"`
	} `json:"datastoresInfo"`
	Health struct {
		UpdatedAt   string `json:"updatedAt"`
		State       string `json:"state"`
		StateReason string `json:"stateReason"`
		Status      string `json:"status"`
	} `json:"health"`
	Network struct {
		DNS []struct {
			NetworkAddress string `json:"networkAddress"`
		} `json:"dns"`
		Hostname string `json:"hostname"`
		Nics     []struct {
			ID             string `json:"id"`
			Name           string `json:"name"`
			NetworkType    string `json:"networkType"`
			NetworkAddress string `json:"networkAddress"`
			NetworkIndex   int    `json:"networkIndex"`
			SubnetMask     string `json:"subnetMask"`
			Gateway        string `json:"gateway"`
		} `json:"nics"`
		Ntp []struct {
			NetworkAddress string `json:"networkAddress"`
		} `json:"ntp"`
		Proxy struct {
			NetworkAddress string `json:"networkAddress"`
			Port           int    `json:"port"`
			Credentials    struct {
				Username string `json:"username"`
			} `json:"credentials"`
		} `json:"proxy"`
	} `json:"network"`
	DataOrchestratorID    string `json:"dataOrchestratorId"`
	SerialNumber          string `json:"serialNumber"`
	SoftwareVersion       string `json:"softwareVersion"`
	VMID                  string `json:"vmId"`
	SupportUserCiphertext string `json:"supportUserCiphertext"`
	AdminUserCiphertext   string `json:"adminUserCiphertext"`
	RemoteAccessEnabled   bool   `json:"remoteAccessEnabled"`
	RemoteAccessStationID string `json:"remoteAccessStationId"`
	Size                  struct {
		MaxOnPremDailyProtectedDataTiB  int `json:"maxOnPremDailyProtectedDataTiB"`
		MaxInCloudDailyProtectedDataTiB int `json:"maxInCloudDailyProtectedDataTiB"`
		MaxOnPremRetentionDays          int `json:"maxOnPremRetentionDays"`
		MaxInCloudRetentionDays         int `json:"maxInCloudRetentionDays"`
	} `json:"size"`
	Override struct {
		CPU        int `json:"cpu"`
		RAMGiB     int `json:"ramGiB"`
		StorageTiB int `json:"storageTiB"`
	} `json:"override"`
	State string `json:"state"`
}
