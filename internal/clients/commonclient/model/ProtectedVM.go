package model

type ProtectedVM struct {
	Identifier           string `json:"identifier"`
	Name                 string `json:"name"`
	ProvisionedStorageMb int    `json:"provisionedStorageMb"`
	UsedStorageMb        int    `json:"usedStorageMb"`
	Vpgs                 []struct {
		Name          string `json:"name"`
		Identifier    string `json:"identifier"`
		Status        string `json:"status"`
		State         string `json:"state"`
		ProtectedSite struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Role string `json:"role"`
		} `json:"protectedSite"`
		RecoverySite struct {
			Name string `json:"name"`
			Type string `json:"type"`
			Role string `json:"role"`
		} `json:"recoverySite"`
	} `json:"vpgs"`
	Zorg struct {
		Name       string `json:"name"`
		Identifier string `json:"identifier"`
	} `json:"zorg"`
}
