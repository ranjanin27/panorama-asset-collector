package model

type MsSqlProtectionGroups struct {
	Items      []MsSqlProtectionGroup `json:"items"`
	PageLimit  int                    `json:"pageLimit"`
	PageOffset int                    `json:"pageOffset"`
	Total      int                    `json:"total"`
}
