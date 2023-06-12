package model

type ProtectionPolicies struct {
	Items      []ProtectionPolicy `json:"items"`
	PageLimit  int                `json:"pageLimit"`
	PageOffset int                `json:"pageOffset"`
	Total      int                `json:"total"`
}
