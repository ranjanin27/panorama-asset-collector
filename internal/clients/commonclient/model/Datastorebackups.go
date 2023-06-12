package model

type DatastoreBackups struct {
	Items      []DatastoreBackup `json:"items"`
	PageLimit  int               `json:"pageLimit"`
	PageOffset int               `json:"pageOffset"`
	Total      int               `json:"total"`
}
