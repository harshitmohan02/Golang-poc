package model

type Totalretained struct {
	TotalActiveResignation int                `json:"totalactiveresignation"`
	TotalPip               int                `json:"totalperformanceimproplan"`
	TotalTbr               int                `json:"totaltoberetained"`
	Data                   []Toberetaineddata `json:"retaineddata"`
}
type Toberetaineddata struct {
	ID                   int    `json:"id,omitempty"`
	ManagerName          string `json:"managername"`
	ProjectName          string `json:"projectname"`
	ActiveResignation    int    `json:"activeresignation"`
	PerformanceImproPlan int    `json:"performanceimproplan"`
	ToBeRetained         int    `json:"toberetained"`
	IsActive             int    `json:"isactive,omitempty"`
}

// type Toberetainedgetall struct {
// 	Manager              string `json:"manager"`
// 	ActiveResignation    int    `json:"activeresignation"`
// 	PerformanceImproPlan int    `json:"performanceimproplan"`
// 	ToBeRetained         int    `json:"toberetained"`
// }
