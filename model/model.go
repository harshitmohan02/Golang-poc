package model

type Pagination struct {
	TotalData  int         `json:"total_data"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
	Page       int         `json:"current_page"`
	Data       interface{} `json:"data"`
}
type UserRole struct {
	ProjectManager int `json:"project_manager"`
	ProgramManager int `json:"program_manager"`
}

type Project struct {
	Id             int    `json:"Id,omitempty"`
	ProjectName    string `json:"project_name,omitempty"`
	SubProjectName string `json:"subproject_name,omitempty"`
	ManagerName    string `json:"project_manager_name,omitempty"`
	ManagerEmailID string `json:"project_manager_email_id,omitempty"`
}

type Position struct {
	ProjectName    string `json:"projectname"`
	Between0to15   int    `json:"between0to15"`
	Between15to30  int    `json:"between15to30"`
	Between30to60  int    `json:"between30to60"`
	Between60to90  int    `json:"between60to90"`
	Between90to120 int    `json:"between90to120"`
	Greaterthen120 int    `json:"greaterthen120"`
}

type Profile struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	ImagePath string `json: "image_path"`
}

type Pages_Daily struct {
	Total_data float64 `json:"total_pages"`
	Limit      float64 `json:"limit"`
	Page       float64 `json:"page"`
	Data       []Daily `json:"posts"`
}

type Pages_Weekly struct {
	Total_data float64  `json:"total_pages"`
	Limit      float64  `json:"limit"`
	Page       float64  `json:"page"`
	Data       []Weekly `json:"posts"`
}

type Daily struct {
	Id                string `json:"id"`
	Project_name      string `json:"project_name"`
	Designation       string `json:"designation"`
	Ageing            int    `json:"ageing`
	Type_position     string `json:"type_position"`
	Position          string `json:"position"`
	Priority          string `json:"priority"`
	Additonal_comment string `json:"additional_comment"`
	L1_due            string `json:"l1_due"`
	L2_due            string `json:"l2_due"`
	Client_due        string `json:"client_due"`
	Is_active         string `json:"is_active"`
	Created_at        string `json:"created_at"`
}

type Weekly struct {
	Id                string `json:"id"`
	Project_name      string `json:"project_name"`
	Designation       string `json:"designation"`
	Ageing            int    `json:"ageing"`
	Type_position     string `json:"type_position"`
	Position          string `json:"position"`
	Priority          string `json:"priority"`
	Additonal_comment string `json:"additional_comment"`
	L1_Happened       string `json:"l1_happened"`
	L2_Happened       string `json:"l2_happened"`
	Client_Happened   string `json:"client_happened"`
	Is_active         string `json:"is_active"`
}
