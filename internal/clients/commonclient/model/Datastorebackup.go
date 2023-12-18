package model

type DatastoreBackup struct {
	AppType        string `json:"appType"`
	BackupMode     string `json:"backupMode"`
	BackupSetsInfo []struct {
		Backups []struct {
			ID             string `json:"id"`
			ObjectCount    int    `json:"objectCount"`
			SizeInBytes    int    `json:"sizeInBytes"`
			SourceCopyInfo []struct {
				CopyType       string `json:"copyType"`
				ID             string `json:"id"`
				Name           string `json:"name"`
				ResourceURI    string `json:"resourceUri"`
				ScsiIdentifier string `json:"scsiIdentifier"`
			} `json:"sourceCopyInfo"`
			SourceVolumeInfo struct {
				ID             string `json:"id"`
				Name           string `json:"name"`
				ResourceURI    string `json:"resourceUri"`
				ScsiIdentifier string `json:"scsiIdentifier"`
			} `json:"sourceVolumeInfo"`
		} `json:"backups"`
		SourceStorageSystemInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"sourceStorageSystemInfo"`
	} `json:"backupSetsInfo"`
	BackupType       string `json:"backupType"`
	Consistency      string `json:"consistency"`
	ContainsRdmDisks bool   `json:"containsRdmDisks"`
	CreatedByInfo    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createdByInfo"`
	DataOrchestratorID  string `json:"dataOrchestratorId"`
	Description         string `json:"description"`
	ExpiresAt           string `json:"expiresAt"`
	Generation          int    `json:"generation"`
	ID                  string `json:"id"`
	Name                string `json:"name"`
	PointInTime         string `json:"pointInTime"`
	ProtectionStoreInfo struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"protectionStoreInfo"`
	ResourceURI  string `json:"resourceUri"`
	ScheduleInfo struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Recurrence string `json:"recurrence"`
	} `json:"scheduleInfo"`
	SourceCopyInfo struct {
		ID          string `json:"id"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"sourceCopyInfo"`
	State             string `json:"state"`
	StateReason       string `json:"stateReason"`
	Status            string `json:"status"`
	StorageSystemInfo struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"storageSystemInfo"`
	Type      string `json:"type"`
	SourceID  string `json:"sourceID"`
	UpdatedAt string `json:"updatedAt"`
}
