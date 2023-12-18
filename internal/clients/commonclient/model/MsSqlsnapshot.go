package model

type MsSqlDBSnapshot struct {
	AssociatedDatabases []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"associatedDatabases"`
	SnapshotType  string `json:"snapshotType"`
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
			DisplayName string `json:"displayName"`
		} `json:"storageSystemInfo"`
	} `json:"logBackupInfo"`
	Name         string `json:"name"`
	PointInTime  string `json:"pointInTime"`
	ResourceURI  string `json:"resourceUri"`
	ScheduleInfo struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Recurrence string `json:"recurrence"`
	} `json:"scheduleInfo"`
	SourceID          string
	State             string `json:"state"`
	StateReason       string `json:"stateReason"`
	Status            string `json:"status"`
	StorageSystemInfo struct {
		DisplayName string `json:"displayName"`
		ID          string `json:"id"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"storageSystemInfo"`
	Type                string `json:"type"`
	UpdatedAt           string `json:"updatedAt"`
	VolumesSnapshotInfo []struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		ResourceURI       string `json:"resourceUri"`
		SCSIIdentifier    string `json:"scsciIdentifier"`
		SizeInMib         int    `json:"sizeInMib"`
		StorageSystemInfo struct {
			ID                string `json:"id"`
			Name              string `json:"name"`
			ResourceURI       string `json:"resourceUri"`
			Type              string `json:"type"`
			DisplayName       string `json:"displayName"`
			Managed           bool   `json:"managed"`
			SerialNumber      string `json:"serialNumber"`
			StorageSystemType string `json:"storageSystemType"`
			VendorName        string `json:"vendorName"`
		} `json:"storageSystemInfo"`
	} `json:"volumesSnapshotInfo"`
}
