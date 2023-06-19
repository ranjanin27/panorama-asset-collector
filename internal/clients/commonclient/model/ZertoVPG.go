package model

type ZertoVPG struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	Name          string `json:"name"`
	ProtectedSite struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"protectedSite"`
	RecoverySite struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"recoverySite"`
	Hypervisor struct {
		Host             string `json:"host"`
		Datastore        string `json:"datastore"`
		Cluster          string `json:"cluster"`
		DatastoreCluster string `json:"datastoreCluster"`
		Folder           string `json:"folder"`
		Network          string `json:"network"`
		TestNetwork      string `json:"testNetwork"`
	} `json:"hypervisor"`
	VirtualMachines []struct {
		ID string `json:"id"`
	} `json:"virtualMachines"`
	Type        string `json:"type"`
	ResourceURI string `json:"resourceUri"`
	CustomerID  string `json:"customerId"`
	Generation  int    `json:"generation"`
	UpdatedAt   string `json:"updatedAt"`
	CreatedAt   string `json:"createdAt"`
}
