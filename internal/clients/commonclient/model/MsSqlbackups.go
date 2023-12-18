package model

type MsSqlDBBackups struct {
	Items      []MsSqlDBBackup `json:"items"`
	PageLimit  int             `json:"pageLimit"`
	PageOffset int             `json:"pageOffset"`
	Total      int             `json:"total"`
}
