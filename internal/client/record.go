package client

import (
	"bufio"
	"schnegge/internal/base"
	"schnegge/internal/config"
	"strconv"
)

func NewDailyReportData(cfg config.Config, orderEmployeeId int) base.NewDailyReportData {
	m, _ := cfg.GetValue(config.Minutes)
	minutes, _ := strconv.Atoi(m)
	h, _ := cfg.GetValue(config.Hours)
	hours, _ := strconv.Atoi(h)
	date, _ := cfg.GetValue(config.Date)
	comment, _ := cfg.GetValue(config.Comment)
	t, _ := cfg.GetValue(config.Training)
	training := t == "true"
	report := base.NewDailyReportData{
		Minutes:         minutes,
		Hours:           hours,
		EmployeeorderId: orderEmployeeId,
		Date:            FormatDateForSalat(CalculateDate(date)[0]),
		Comment:         comment,
		Training:        training,
	}
	return report
}

func Record(cfg config.Config, buchung base.NewDailyReportData) {
	if !checkConfig(cfg) {
		return
	}
	resp, err := doPostRequest(cfg, "/rest/daily-reports/", buchung)

	if err != nil {
		base.Log.Panic(err)
	}
	defer resp.Body.Close()

	base.Log.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 50; i++ {
		base.Log.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		base.Log.Panic(err)
	}

	if resp.StatusCode != 201 {
		base.Log.Panic("Neue DailyReportData ist fehlgeschlagen. ", resp.Status, buchung)
	}
}
