package model

type Datastore struct {
	AllowedOperations   []string `json:"allowedOperations"`
	AppType             string   `json:"appType"`
	CapacityFree        int64    `json:"capacityFree"`
	CapacityInBytes     int64    `json:"capacityInBytes"`
	CapacityUncommitted int      `json:"capacityUncommitted"`
	ClusterInfo         struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"clusterInfo"`
	CreatedAt       string `json:"createdAt"`
	CustomerID      string `json:"customerId"`
	DatacentersInfo []struct {
		ID    string `json:"id"`
		Moref string `json:"moref"`
		Name  string `json:"name"`
	} `json:"datacentersInfo"`
	DatastoreType string `json:"datastoreType"`
	DisplayName   string `json:"displayName"`
	FolderInfo    struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"folderInfo"`
	Generation     int    `json:"generation"`
	HciClusterUUID string `json:"hciClusterUuid"`
	HostsInfo      []struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"hostsInfo"`
	HypervisorManagerInfo struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"hypervisorManagerInfo"`
	ID                string `json:"id"`
	Moref             string `json:"moref"`
	Name              string `json:"name"`
	Protected         bool   `json:"protected"`
	ProtectionJobInfo struct {
		ID                   string `json:"id"`
		ProtectionPolicyInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"protectionPolicyInfo"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"protectionJobInfo"`
	ProtectionPolicyAppliedAtInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"protectionPolicyAppliedAtInfo"`
	ProvisioningPolicyInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"provisioningPolicyInfo"`
	ReplicationInfo struct {
		ID             string      `json:"id"`
		Name           string      `json:"name"`
		PartnerDetails interface{} `json:"partnerDetails"`
		ResourceURI    string      `json:"resourceUri"`
	} `json:"replicationInfo"`
	ResourceURI            string      `json:"resourceUri"`
	Services               []string    `json:"services"`
	State                  string      `json:"state"`
	StateReason            string      `json:"stateReason"`
	Status                 string      `json:"status"`
	Type                   string      `json:"type"`
	UID                    string      `json:"uid"`
	UpdatedAt              string      `json:"updatedAt"`
	VMCount                int         `json:"vmCount"`
	VMProtectionGroupsInfo interface{} `json:"vmProtectionGroupsInfo"`
	VolumesInfo            interface{} `json:"volumesInfo"`
}
