package model

type CSPAccount struct {
	CustomerID         string      `json:"customerId"`
	Generation         int         `json:"generation"`
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	ConsoleURI         string      `json:"consoleUri"`
	ResourceURI        string      `json:"resourceUri"`
	Type               string      `json:"type"`
	Suspended          bool        `json:"suspended"`
	CspType            string      `json:"cspType"`
	CspID              string      `json:"cspId"`
	ValidatedAt        string      `json:"validatedAt"`
	ValidationErrors   interface{} `json:"validationErrors"`
	ValidationStatus   string      `json:"validationStatus"`
	RefreshStatus      string      `json:"refreshStatus"`
	RefreshedAt        string      `json:"refreshedAt"`
	OnboardingTemplate struct {
		VersionApplied string `json:"versionApplied"`
		LatestVersion  string `json:"latestVersion"`
		UpgradeNeeded  bool   `json:"upgradeNeeded"`
		Message        string `json:"message"`
	} `json:"onboardingTemplate"`
}
