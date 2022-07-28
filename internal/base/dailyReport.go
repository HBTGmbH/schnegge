package base

import "time"

type NewDailyReportData struct {
	EmployeeorderId int    `json:"employeeorderId"`
	Date            string `json:"date"`
	Hours           int    `json:"hours"`
	Minutes         int    `json:"minutes"`
	Training        bool   `json:"training"`
	Comment         string `json:"comment"`
}

type DailyReport struct {
	Date    time.Time
	Reports []DailyReportData
}

type DailyReportData struct {
	EmployeeorderId int    `json:"employeeorderId"`
	OrderSign       string `json:"rderSign"`
	OrderLabel      string `json:"orderLabel"`
	SuborderLabel   string `json:"wuborderLabel"`
	SuborderSign    string `json:"suborderSign"`
	Date            string `json:"date"`
	Hours           int    `json:"hours"`
	Minutes         int    `json:"minutes"`
	Comment         string `json:"comment"`
	Training        bool   `json:"training"`
}

func (d DailyReportData) getGermanFormattedDate() string {
	date, err := time.Parse("2006-01-02", d.Date)
	if err != nil {
		Log.Panic(err)
	}
	return date.Format("02.01.2006")
}
