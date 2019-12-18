package model

type Pagenation struct{
    TotalData int `json:"total_data"`
    Limit     int  `json:"limit"`
    TotalPages int `json:"total_pages"`
    Page      int `json:"current_page"`
    Data    []Project `json:"data"`

}

type Project struct {
    Id             int `json:"Id,omitempty"`
    ProjectName    string `json:"name,omitempty"`
    ManagerName    string `json:"manager_name,omitempty"`
    ManagerEmailID string `json:"manager_email_id,omitempty"`
    Flag           string `json:"flag,omitempty"`
}

type Position struct {
    ProjectName    string `json:"projectname"` 
    Between0to15   int `json:"between0to15"`
    Between15to30  int `json:"between15to30"`
    Between30to60  int `json:"between30to60"`
    Between60to90  int `json:"between60to90"`
    Between90to120 int `json:"between90to120"`
    Greaterthen120 int `json:"greaterthen120"`
}

type Conf struct{
  
    R []Routes `yaml:"Routes"`
}
type Routes struct{
    Path string `yaml:"Path"`
    Callback string `yaml:"Callback"`
    Method string `yaml:"Method"`
    Authorization string `yaml:"Authorization"`
}
