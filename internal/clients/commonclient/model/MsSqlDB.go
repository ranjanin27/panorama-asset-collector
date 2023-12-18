package model

type MsSqlDB struct {
	ApplicationHostInfo struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"applicationHostInfo"`
	AvailabilityGroupInfo struct {
		Name string `json:"name"`
		Role string `json:"role"`
		Uid  string `json:"uid"`
	} `json:"availabilityGroupInfo"`
	ClusterInfo struct {
		Clustered bool   `json:"clustered"`
		Role      string `json:"role"`
	} `json:"clusterInfo"`
	CreatedAt    string `json:"createdAt"`
	CustomerID   string `json:"customerId"`
	Generation   int    `json:"generation"`
	ID           string `json:"id"`
	InstanceInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"instanceInfo"`
	MssqlDatabaseProtectionGroupInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"mssqlDatabaseProtectionGroupInfo"`
	Name              string `json:"name"`
	ProtectionJobInfo struct {
		ID                   string `json:"id"`
		Name                 string `json:"name"`
		ResourceURI          string `json:"resourceUri"`
		Type                 string `json:"type"`
		ProtectionPolicyInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"protectionPolicyInfo"`
	} `json:"protectionJobInfo"`
	ProtectionStatus    string `json:"protectionStatus"`
	RecoveryPointsExist bool   `json:"recoveryPointsExist"`
	ResourceURI         string `json:"resourceUri"`
	SizeInBytes         string `json:"sizeInBytes"`
	State               string `json:"state"`
	StateReason         string `json:"stateReason"`
	Status              string `json:"status"`
	Type                string `json:"type"`
	UpdatedAt           string `json:"updatedAt"`
	VirtualizationInfo  struct {
		HypervisorManagerInfo struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"hypervisorManagerInfo"`
		VirtualMachineInfo struct {
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"virtualMachineInfo"`
	} `json:"virtualizationInfo"`
}
