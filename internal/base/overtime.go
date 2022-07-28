package base

type Overtime struct {
	Total        Total        `json:"total"`
	CurrentMonth CurrentMonth `json:"currentMonth"`
}

type Total struct {
	Begin    string `json:"begin"`
	End      string `json:"end"`
	Days     int    `json:"days"`
	Hours    int    `json:"hours"`
	Minutes  int    `json:"minutes"`
	Duration string `json:"duration"`
	Negative bool   `json:"negative"`
}

type CurrentMonth struct {
	Begin    string `json:"begin"`
	End      string `json:"end"`
	Days     int    `json:"days"`
	Hours    int    `json:"hours"`
	Minutes  int    `json:"minutes"`
	Duration string `json:"duration"`
	Negative bool   `json:"negative"`
}
