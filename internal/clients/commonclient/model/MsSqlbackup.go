package model

type MsSqlDBBackup struct {
	AssociatedDatabases []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"associatedDatabases"`
	BackupGranularity string `json:"backupGranularity"`
	BackupSetsInfo    []struct {
		Backups []struct {
			ID          string `json:"id"`
			ObjectCount int    `json:"objectCount"`
			SizeInBytes int    `json:"sizeInBytes"`
		} `json:"backups"`
		SourceStorageSystemInfo struct {
			DisplayName string `json:"displayName"`
			ID          string `json:"id"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
		} `json:"sourceStorageSystemInfo"`
	} `json:"backupSetsInfo"`
	BackupType    string `json:"backupType"`
	CreatedByInfo struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createdByInfo"`
	CustomerID    string `json:"customerId"`
	Description   string `json:"description"`
	ExpiresAt     string `json:"expiresAt"`
	Generation    int    `json:"generation"`
	ID            string `json:"id"`
	LockedUntil   string `json:"lockedUntil"`
	LogBackupInfo struct {
		LastLogBackupTime   string `json:"lastLogBackupTime"`
		ProtectionStoreInfo struct {
			ID                  string `json:"id"`
			Name                string `json:"name"`
			ResourceURI         string `json:"resourceUri"`
			Type                string `json:"type"`
			ProtectionStoreType string `json:"protectionStoreType"`
		} `json:"protectionStoreInfo"`
		StorageSystemInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			ResourceURI string `json:"resourceUri"`
			Type        string `json:"type"`
			DisplayName string `json:"displayNmae"`
		} `json:"storageSystemInfo"`
	} `json:"logBackupInfo"`
	Name                string `json:"name"`
	PointInTime         string `json:"pointInTime"`
	ProtectionStoreInfo struct {
		ID                  string `json:"id"`
		Name                string `json:"name"`
		ResourceURI         string `json:"resourceUri"`
		Type                string `json:"type"`
		ProtectionStoreType string `json:"protectionStoreType"`
	} `json:"protectionStoreInfo"`
	ResourceURI  string `json:"resourceUri"`
	ScheduleInfo struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Recurrence string `json:"recurrence"`
	} `json:"scheduleInfo"`
	SourceID       string
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
	UpdatedAt string `json:"updatedAt"`
}
