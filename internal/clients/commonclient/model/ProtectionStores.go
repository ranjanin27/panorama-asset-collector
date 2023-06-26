package model

type ProtectionStores struct {
	Items      []ProtectionStore `json:"items"`
	PageLimit  int               `json:"pageLimit"`
	PageOffset int               `json:"pageOffset"`
	Total      int               `json:"total"`
}
