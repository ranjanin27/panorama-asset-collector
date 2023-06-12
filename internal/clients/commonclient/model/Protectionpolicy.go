package model

type ProtectionPolicy struct {
	Assigned           bool `json:"assigned"`
	ProtectionJobsInfo []struct {
		ID        string `json:"id"`
		AssetInfo struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
			Type        string `json:"type"`
			ResourceURI string `json:"resourceUri"`
		} `json:"assetInfo"`
		ResourceURI string `json:"resourceUri"`
		Type        string `json:"type"`
	} `json:"protectionJobsInfo,omitempty"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Protections []struct {
		ID        string `json:"id"`
		Schedules []struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			Schedule struct {
				Recurrence     string `json:"recurrence"`
				RepeatInterval struct {
					Every int `json:"every"`
				} `json:"repeatInterval"`
				StartTime string `json:"startTime"`
			} `json:"schedule"`
			ExpireAfter struct {
				Unit  string `json:"unit"`
				Value int    `json:"value"`
			} `json:"expireAfter"`
			NamePattern struct {
				Format string `json:"format"`
			} `json:"namePattern"`
		} `json:"schedules"`
		ProtectionStoreInfo struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Type        string `json:"type"`
			ResourceURI string `json:"resourceUri"`
		} `json:"protectionStoreInfo"`
		Type            string `json:"type"`
		ApplicationType string `json:"applicationType"`
	} `json:"protections"`
	CreatedAt string `json:"createdAt"`
	CreatedBy struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"createdBy"`
	Generation  int    `json:"generation"`
	ResourceURI string `json:"resourceUri"`
	ConsoleURI  string `json:"consoleUri"`
	Type        string `json:"type"`
	UpdatedAt   string `json:"updatedAt"`
}
