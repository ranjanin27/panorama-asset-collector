package model

type ZertoVPGs struct {
	Items  []ZertoVPG `json:"items"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
	Total  int        `json:"total"`
}
