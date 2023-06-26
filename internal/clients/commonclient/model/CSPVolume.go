package model

type CSPVolume struct {
	CustomerID  string `json:"customerId"`
	Generation  int    `json:"generation"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	ResourceURI string `json:"resourceUri"`
	ConsoleURI  string `json:"consoleUri"`
	Type        string `json:"type"`
	AccountID   string `json:"accountId"`
	CspInfo     struct {
		AvailabilityZone string        `json:"availabilityZone"`
		CreatedAt        string        `json:"createdAt"`
		ID               string        `json:"id"`
		IsEncrypted      bool          `json:"isEncrypted"`
		Region           string        `json:"region"`
		SizeInGiB        int           `json:"sizeInGiB"`
		Iops             int           `json:"iops"`
		Tags             []interface{} `json:"tags"`
		VolumeType       string        `json:"volumeType"`
	} `json:"cspInfo"`
	MachineInstanceAttachmentInfo []struct {
		AttachedTo struct {
			Type        string `json:"type"`
			ResourceURI string `json:"resourceUri"`
			Name        string `json:"name"`
		} `json:"attachedTo"`
		State               string `json:"state"`
		Device              string `json:"device"`
		DeleteOnTermination bool   `json:"deleteOnTermination"`
		AttachedAt          string `json:"attachedAt"`
	} `json:"machineInstanceAttachmentInfo"`
	State               string        `json:"state"`
	ProtectionGroupInfo []interface{} `json:"protectionGroupInfo"`
	ProtectionStatus    string        `json:"protectionStatus"`
	ProtectionJobInfo   []interface{} `json:"protectionJobInfo"`
	BackupInfo          []interface{} `json:"backupInfo"`
}
