package model

type VirtualMachines struct {
	Items      []VirtualMachine `json:"items"`
	PageLimit  int              `json:"pageLimit"`
	PageOffset int              `json:"pageOffset"`
	Total      int              `json:"total"`
}
