package model

type Resignations struct { //Resignations structure
	ID                int    `json:"id,omitempty"`
	Empname           string `json:"empname,omitempty"`
	Project           string `json:"project,omitempty"`
	Manager           string `json:"manager,omitempty"`
	Backfillrequired  string `json:"backfillrequired,omitempty"`
	Regrenonregre     string `json:"regrenonregre,omitempty"`
	Status            string `json:"status,omitempty"`
	Dateofresignation string `json:"dor,omitempty"`
	Dateofleaving     string `json:"dol,omitempty"`
	CreatedAt         string `json:"createdat,omitempty"`
	UpdatedAt         string `json:"updatedat,omitempty"`
	IsActive          int    `json:"isactive,omitempty"`
}
