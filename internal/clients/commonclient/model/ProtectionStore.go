package model

type ProtectionStore struct {
	Name              string `json:"name"`
	DisplayName       string `json:"displayName"`
	Description       string `json:"description"`
	ResourceURI       string `json:"resourceUri"`
	ID                string `json:"id"`
	Status            string `json:"status"`
	State             string `json:"state"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	StorageSystemInfo struct {
		Type        string `json:"type"`
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		ResourceURI string `json:"resourceUri"`
	} `json:"storageSystemInfo"`
	SizeOnDiskInBytes     int    `json:"sizeOnDiskInBytes"`
	UserDataStoredInBytes int64  `json:"userDataStoredInBytes"`
	MaxCapacityInBytes    int    `json:"maxCapacityInBytes"`
	ProtectionStoreType   string `json:"protectionStoreType"`
	Type                  string `json:"type"`
	Region                string `json:"region,omitempty"`
}
