package model

type Storeonces struct {
	Items      []Storeonce `json:"items"`
	PageLimit  int         `json:"pageLimit"`
	PageOffset int         `json:"pageOffset"`
	Total      int         `json:"total"`
}
