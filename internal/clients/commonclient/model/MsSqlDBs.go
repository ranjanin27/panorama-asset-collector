package model

type MsSqlDBs struct {
	Items      []MsSqlDB `json:"items"`
	PageLimit  int       `json:"pageLimit"`
	PageOffset int       `json:"pageOffset"`
	Total      int       `json:"total"`
}
