package model

type Storeonce struct {
	CustomerID  string `json:"customerId"`
	Generation  int    `json:"generation"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	ResourceURI string `json:"resourceUri"`
	DateTime    struct {
		MethodDateTimeSet string `json:"methodDateTimeSet"`
		Timezone          string `json:"timezone"`
		UtcDateTime       string `json:"utcDateTime"`
	} `json:"dateTime"`
	Health struct {
		State       string `json:"state"`
		StateReason string `json:"stateReason"`
		Status      string `json:"status"`
		UpdatedAt   string `json:"updatedAt"`
	} `json:"health"`
	Network struct {
		DNS      []interface{} `json:"dns"`
		Hostname string        `json:"hostname"`
		Ntp      []interface{} `json:"ntp"`
		Nics     []interface{} `json:"nics"`
	} `json:"network"`
	SerialNumber    string `json:"serialNumber"`
	SoftwareVersion string `json:"softwareVersion"`
	Storage         struct {
		State                    string `json:"state"`
		UnconfiguredStorageBytes int    `json:"unconfiguredStorageBytes"`
		ConfiguredStorageBytes   int    `json:"configuredStorageBytes"`
		CapacityLicensedBytes    int    `json:"capacityLicensedBytes"`
		CapacityUnlicensedBytes  int    `json:"capacityUnlicensedBytes"`
		UsedBytes                int    `json:"usedBytes"`
		FreeBytes                int    `json:"freeBytes"`
	} `json:"storage"`
	FibreChannel struct {
		Initiators []interface{} `json:"initiators"`
	} `json:"fibreChannel"`
	ISCSIInitiatorName string        `json:"iSCSIInitiatorName"`
	DataOrchestrators  []interface{} `json:"dataOrchestrators"`
	Type               string        `json:"type"`
	Description        string        `json:"description"`
}
