package model

type Active struct { //Structure of total active resignation
	ActiveHeadCount   string `json:"activeheadcount"`
	Billable          string `json:"billable"`
	BillingOnHold     string `json:"billingonhold"`
	ValueTrade        string `json:"valuetrade"`
	ProjectInvestment string `json:"projectinvestment"`
	Others            string `json:"others"`
}
