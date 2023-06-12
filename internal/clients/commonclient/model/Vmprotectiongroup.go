package model

type VMProtectionGroup struct {
	AppType               string `json:"appType"`
	CreatedAt             string `json:"createdAt"`
	Generation            int    `json:"generation"`
	HypervisorManagerInfo struct {
		Name        string `json:"name"`
		ID          string `json:"id"`
		ResourceURI string `json:"resourceUri"`
	} `json:"hypervisorManagerInfo"`
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	DataOrchestratorID     string `json:"dataOrchestratorId"`
	ResourceURI            string `json:"resourceUri"`
	ConsoleURI             string `json:"consoleUri"`
	Type                   string `json:"type"`
	UpdatedAt              string `json:"updatedAt"`
	VMProtectionGroupType  string `json:"vmProtectionGroupType"`
	ProtectedResourcesInfo struct {
		VirtualMachinesCount int `json:"virtualMachinesCount"`
	} `json:"protectedResourcesInfo"`
	ProtectionJobInfo struct {
		ProtectionPolicyInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"protectionPolicyInfo"`
		ID          string `json:"id"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"protectionJobInfo"`
	AssetsCategory string `json:"assetsCategory"`
	Assets         []struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		Type        string `json:"type"`
		ResourceURI string `json:"resourceUri"`
	} `json:"assets"`
	State       string `json:"state"`
	StateReason string `json:"stateReason"`
	Status      string `json:"status"`
}
