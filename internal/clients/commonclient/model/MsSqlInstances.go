package model

type MsSqlInstances struct {
	Items      []MsSqlInstance `json:"items"`
	PageLimit  int             `json:"pageLimit"`
	PageOffset int             `json:"pageOffset"`
	Total      int             `json:"total"`
}
