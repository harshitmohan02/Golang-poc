package model

type HeadCount struct {
	ID             int    `json:"id,omitempty"`
	ProjectName    string `json:"projectname"`
	ManagerName    string `json:"managername"`
	BillablesCount int    `json:"billablescount"`
	BillingOnHold  int    `json:"billingonhold"`
	VtCount        int    `json:"vtcount"`
	PiICount       int    `json:"piicount"`
	Others         int    `json:"others"`
	Net            int    `json:"net"`
	CreatedAt      string `json:"createdat,omitempty"`
	UpdatedAt      string `json:"updatedat,omitempty"`
	IsActive       int    `json:"is_active,omitempty"`
}
