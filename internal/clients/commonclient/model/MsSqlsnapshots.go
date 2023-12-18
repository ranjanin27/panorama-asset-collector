package model

type MsSqlDBSnapshots struct {
	Items      []MsSqlDBSnapshot `json:"items"`
	PageLimit  int               `json:"pageLimit"`
	PageOffset int               `json:"pageOffset"`
	Total      int               `json:"total"`
}
