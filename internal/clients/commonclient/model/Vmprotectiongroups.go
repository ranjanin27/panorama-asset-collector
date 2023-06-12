package model

type VMProtectionGroups struct {
	Items      []VMProtectionGroup `json:"items"`
	PageLimit  int                 `json:"pageLimit"`
	PageOffset int                 `json:"pageOffset"`
	Total      int                 `json:"total"`
}
