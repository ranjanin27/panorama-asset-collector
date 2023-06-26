package model

type CSPAccounts struct {
	Items      []CSPAccount `json:"items"`
	PageLimit  int          `json:"pageLimit"`
	PageOffset int          `json:"pageOffset"`
	Total      int          `json:"total"`
}
