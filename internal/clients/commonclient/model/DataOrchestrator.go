package model

type DO struct {
	ConnectionState string `json:"connectionState"`
	ConsoleUri      string `json:"consoleUri"`
	CustomerID      string `json:"customerId"`
	DisplayName     string `json:"displayName"`
	Generation      int    `json:"generation"`
	ID              string `json:"id"`
	Name            string `json:"name"`
	ResourceURI     string `json:"resourceUri"`
	SerialNumber    string `json:"serialNumber"`
	SoftwareVersion string `json:"softwareVersion"`
	State           string `json:"state"`
	StateReason     string `json:"stateReason"`
	Status          string `json:"status"`
	Type            string `json:"type"`
	UpTimeInSeconds string `json:"upTimeInSeconds"`
	CreatedAt       string `json:"createdAt"`
	DateTime        struct {
		MethodDateTimeSet string `json:"methodDateTimeSet"`
		Timezone          string `json:"timezone"`
		UtcDateTime       string `json:"utcDateTime"`
	} `json:"dateTime"`
	Interfaces struct {
		Network interface{} `json:"network"`
	}
	LastUpdateCheckTime string `json:"lastUpdateCheckTime"`
	LatestRecoveryPoint string `json:"latestRecoveryPoint"`
	NTP                 struct {
		NtpServers  []interface{} `json:"ntpServers"`
		State       string        `json:"state"`
		StateReason string        `json:"stateReason"`
		Status      string        `json:"status"`
	}
	Platform              string `json:"platform"`
	PoweredOnAt           string `json:"poweredOnAt"`
	TotalMemoryInGiB      int    `json:"totalMemoryInGiB"`
	UpdatedAt             string `json:"updatedAt"`
	VCpu                  int    `json:"vCpu"`
	AdminUserCiphertext   string `json:"adminUserCiphertext"`
	InfosightEnabled      bool   `json:"infosightEnabled"`
	RemoteAccessEnabled   bool   `json:"remoteAccessEnabled"`
	RemoteAccessStationId string `json:"remoteAccessStationId"`
	SupportUserCiphertext string `json:"supportUserCiphertext"`
}
