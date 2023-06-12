package model

type DSSnapshots struct {
	Items      []DSSnapshot `json:"items"`
	PageLimit  int          `json:"pageLimit"`
	PageOffset int          `json:"pageOffset"`
	Total      int          `json:"total"`
}
