package model

type Datastores struct {
	Items      []Datastore `json:"items"`
	PageLimit  int         `json:"pageLimit"`
	PageOffset int         `json:"pageOffset"`
	Total      int         `json:"total"`
}
