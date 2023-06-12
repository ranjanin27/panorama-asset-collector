package model

type ProtectedVMs struct {
	Items      []ProtectedVM `json:"items"`
	PageLimit  int           `json:"pageLimit"`
	PageOffset int           `json:"pageOffset"`
	Total      int           `json:"total"`
}
