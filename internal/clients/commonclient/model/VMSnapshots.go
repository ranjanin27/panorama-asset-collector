package model

type VMSnapshots struct {
	Items      []VMSnapshot `json:"items"`
	PageLimit  int          `json:"pageLimit"`
	PageOffset int          `json:"pageOffset"`
	Total      int          `json:"total"`
}
