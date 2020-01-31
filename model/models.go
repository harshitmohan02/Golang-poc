package model

type Overview struct {
	GrandTotal        int64 `json:"grandtotal"`
	ClosedInTime      int64 `json:"closedintime"`
	ClosedOutTime     int64 `json:"closedouttime"`
	InProgressInTime  int64 `json:"inprogressintime"`
	InProgressOutTime int64 `json:"inprogressouttime"`
	OnHold            int64 `json:"onhold"`
}

type Manager struct {
	Name  string  `json:"name"`
	SayDo float64 `json:"saydo"`
}

type Saydo struct {
	WeekStart string  `json:"weekstart"`
	WeekEnd   string  `json:"weekend"`
	Saydo     float32 `json:"saydo"`
}

type Prospects struct {
	ID         string `json:"id"`
	Project    string `json:"project"`
	Manager    string `json:"manager"`
	Prospect   string `json:"prospect"`
	Status     string `json:"status"`
	Comments   string `json:"comments"`
	Challenges string `json:"challenges"`
}

type Updates struct {
	ID             string `json:"id"`
	Manager        string `json:"manager"`
	ProjectName    string `json:"projectname"`
	Ups            string `json:"ups"`
	Downs          string `json:"downs"`
	ProjectUpdates string `json:"project_updates"`
	GeneralUpdates string `json:"general_updates"`
	Challenges     string `json:"challenges"`
	NeedHelp       string `json:"need_help"`
	ClientVisits   string `json:"client_visits"`
	TeamSize       string `json:"team_size"`
	OpenPositions  string `json:"open_positions"`
	HighPerformer  string `json:"high_performer"`
	LowPerformer   string `json:"low_performer"`
	// Is_active      string `json:"is_active"`
}
