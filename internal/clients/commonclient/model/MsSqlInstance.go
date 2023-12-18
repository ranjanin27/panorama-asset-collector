package model

type MsSqlInstance struct {
	ApplicationHostInfo struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"applicationHostInfo"`
	AvailabilityGroupsInfo []interface{} `json:"availabilityGroupsInfo"`
	Clustered              bool          `json:"clustered"`
	CreatedAt              string        `json:"createdAt"`
	Credentials            struct {
		Mode     string `json:"mode"`
		Username string `json:"username"`
	} `json:"credentials"`
	CustomerID         string `json:"customerId"`
	Generation         int    `json:"generation"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	ProductName        string `json:"productName"`
	ProductVersion     string `json:"productVersion"`
	ResourceURI        string `json:"resourceUri"`
	State              string `json:"state"`
	StateReason        string `json:"stateReason"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	UpdatedAt          string `json:"updatedAt"`
	VirtualizationInfo struct {
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
