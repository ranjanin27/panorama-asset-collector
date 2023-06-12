package model

type CSPMachineInstance struct {
	CustomerID  string `json:"customerId"`
	Generation  int    `json:"generation"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	ResourceURI string `json:"resourceUri"`
	ConsoleURI  string `json:"consoleUri"`
	Type        string `json:"type"`
	AccountID   string `json:"accountId"`
	CspInfo     struct {
		AccessProfileID  interface{} `json:"accessProfileId"`
		AvailabilityZone string      `json:"availabilityZone"`
		CPUCoreCount     int         `json:"cpuCoreCount"`
		CreatedAt        string      `json:"createdAt"`
		ID               string      `json:"id"`
		InstanceType     string      `json:"instanceType"`
		KeyPairName      string      `json:"keyPairName"`
		NetworkInfo      struct {
			PrivateIPAddress          string      `json:"privateIpAddress"`
			PublicIPAddress           interface{} `json:"publicIpAddress"`
			PublicIPAddressIsFloating bool        `json:"publicIpAddressIsFloating"`
			SecurityGroups            []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"securityGroups"`
			Subnet struct {
				ID string `json:"id"`
			} `json:"subnet"`
			Vpc struct {
				ID string `json:"id"`
			} `json:"vpc"`
		} `json:"networkInfo"`
		Platform   string `json:"platform"`
		Region     string `json:"region"`
		RootDevice string `json:"rootDevice"`
		State      string `json:"state"`
		Tags       []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"tags"`
		VirtualizationType string `json:"virtualizationType"`
	} `json:"cspInfo"`
	VolumeAttachmentInfo []struct {
		AttachedTo struct {
			Type        string `json:"type"`
			ResourceURI string `json:"resourceUri"`
			Name        string `json:"name"`
		} `json:"attachedTo"`
		State               string `json:"state"`
		Device              string `json:"device"`
		DeleteOnTermination bool   `json:"deleteOnTermination"`
		AttachedAt          string `json:"attachedAt"`
	} `json:"volumeAttachmentInfo"`
	State               string        `json:"state"`
	ProtectionGroupInfo []interface{} `json:"protectionGroupInfo"`
	ProtectionStatus    string        `json:"protectionStatus"`
	ProtectionJobInfo   []interface{} `json:"protectionJobInfo"`
	BackupInfo          []interface{} `json:"backupInfo"`
}
