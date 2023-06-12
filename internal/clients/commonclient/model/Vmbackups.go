package model

type VMBackups struct {
	Items      []VMBackup `json:"items"`
	PageLimit  int        `json:"pageLimit"`
	PageOffset int        `json:"pageOffset"`
	Total      int        `json:"total"`
}
