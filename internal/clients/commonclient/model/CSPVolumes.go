package model

type CSPVolumes struct {
	Items      []CSPVolume `json:"items"`
	PageLimit  int         `json:"pageLimit"`
	PageOffset int         `json:"pageOffset"`
	Total      int         `json:"total"`
}
