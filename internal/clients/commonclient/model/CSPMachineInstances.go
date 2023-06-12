package model

type CSPMachineInstances struct {
	Items      []CSPMachineInstance `json:"items"`
	PageLimit  int                  `json:"pageLimit"`
	PageOffset int                  `json:"pageOffset"`
	Total      int                  `json:"total"`
}
