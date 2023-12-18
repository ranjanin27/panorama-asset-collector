package model

type MsSqlProtectionGroup struct {
	CreatedAt   string `json:"createdAt"`
	CustomerID  string `json:"customerId"`
	Description string `json:"description"`
	Generation  int    `json:"generation"`
	ID          string `json:"id"`
	Members     []struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		ResourceURI  string `json:"resourceUri"`
		Type         string `json:"type"`
		InstanceInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"instanceInfo"`
	} `json:"members"`
	Name          string `json:"name"`
	NativeAppInfo struct {
		AvailabilityGroupReplicas []struct {
			Role        string `json:"role"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"availabilityGroupReplicas"`
		ExcludeSystemDatabases bool   `json:"excludeSystemDatabases"`
		ID                     string `json:"id"`
		InstanceInfo           struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"instanceInfo"`
		Name string `json:"name"`
		Type string `json:"type"`
		UID  string `json:"uid"`
	} `json:"nativeAppInfo"`

	ProtectionGroupType string `json:"protectionGroupType"`
	ProtectionJobInfo   struct {
		ID                   string `json:"id"`
		Name                 string `json:"name"`
		ResourceURI          string `json:"resourceUri"`
		Type                 string `json:"type"`
		ProtectionPolicyInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		}
	} `json:"protectionJobInfo"`
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
	} `json:"virtualizationInfo"`
}
