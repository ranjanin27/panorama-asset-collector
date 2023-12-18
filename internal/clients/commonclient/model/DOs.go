package model

type DOs struct {
	Items      []DO `json:"items"`
	PageLimit  int  `json:"pageLimit"`
	PageOffset int  `json:"pageOffset"`
	Total      int  `json:"total"`
}
