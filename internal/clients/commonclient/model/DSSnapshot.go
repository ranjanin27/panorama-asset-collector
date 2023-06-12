package model

type DSSnapshot struct {
	Consistency      string `json:"consistency"`
	ContainsRdmDisks bool   `json:"containsRdmDisks"`
	CreatedByInfo    struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createdByInfo"`
	DataOrchestratorID string `json:"dataOrchestratorId"`
	Description        string `json:"description"`
	ExpiresAt          string `json:"expiresAt"`
	Generation         string `json:"generation"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	PointInTime        string `json:"pointInTime"`
	ResourceURI        string `json:"resourceUri"`
	ScheduleInfo       struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Recurrence string `json:"recurrence"`
	} `json:"scheduleInfo"`
	SnapshotType       string `json:"snapshotType"`
	State              string `json:"state"`
	StateReason        string `json:"stateReason"`
	Status             string `json:"status"`
	StorageSystemsInfo []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"storageSystemsInfo"`
	Type                string `json:"type"`
	UpdatedAt           string `json:"updatedAt"`
	VolumesSnapshotInfo []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		ScsiIdentifier string `json:"scsiIdentifier"`
		SizeInMiB      int    `json:"sizeInMiB"`
	} `json:"volumesSnapshotInfo"`
	PageLimit  int `json:"pageLimit"`
	PageOffset int `json:"pageOffset"`
	Total      int `json:"total"`
}
