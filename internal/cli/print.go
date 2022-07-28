package cli

import (
	"fmt"
	"schnegge/internal/base"
	"sort"
	"time"
)

func PrintSplashScreen() {
	fmt.Println("     .----.   @   @")
	fmt.Println("    / .-\"-.`.  \\v/")
	fmt.Println("    | | '\\ \\ \\_/ )")
	fmt.Println("  ,-\\ `-.' /.'  /")
	fmt.Print(" '---`----'----'")
	fmt.Println("       Kleine Schneggen mögen gerne Salat ...")
}

func PrintOvertime(overtime base.Overtime) {
	fmt.Printf("Überstunden Total:            %3dd, %dh, %02dm \nÜberstunden aktueller Monat:  %3dd, %dh, %02dm\n",
		overtime.Total.Days,
		overtime.Total.Hours,
		overtime.Total.Minutes,
		overtime.CurrentMonth.Days,
		overtime.CurrentMonth.Hours,
		overtime.CurrentMonth.Minutes,
	)
}

func PrintDailyReports(orderLists []base.DailyReport) {
	for _, orderList := range orderLists {
		printDailyReport(orderList)
	}
}

func PrintByProject(orderLists []base.DailyReport) {
	splittedList := split(orderLists)
	sort.SliceStable(splittedList, func(i, j int) bool {
		return compare(splittedList[i], splittedList[j])
	})
	lastElement := len(splittedList) - 1
	startIndex := 0
	subOrderSign := ""
	for i, d := range splittedList {
		r := d.Reports[0]
		if r.SuborderSign != subOrderSign || i == lastElement {
			if i == lastElement {
				printProject(splittedList[startIndex : i+1])
			} else if i != 0 {
				printProject(splittedList[startIndex:i])
			}
			startIndex = i
			subOrderSign = r.SuborderSign
		}
	}
}

func printDailyReport(dailyReport base.DailyReport) {
	printDateHeader(dailyReport)
	for i, order := range dailyReport.Reports {
		printSingleLine(i, order)
	}
}

func printDateHeader(dailyReport base.DailyReport) {
	sumH, sumM := 0, 0
	for _, report := range dailyReport.Reports {
		sumH = sumH + report.Hours
		sumM = sumM + report.Minutes
		if sumM >= 60 {
			sumH = sumH + 1
			sumM = sumM - 60
		}
	}
	fmt.Println(
		base.Bold(
			fmt.Sprintf("=== %v %v (%dh %2dm) ===",
				formatWeekdayGermanShort(dailyReport.Date), formatDateGerman(dailyReport.Date), sumH, sumM)))
}

func printProject(dailyReports []base.DailyReport) {
	hours := 0
	minutes := 0
	for _, dr := range dailyReports {
		hours = hours + dr.Reports[0].Hours
		minutes = minutes + dr.Reports[0].Minutes
	}
	first := dailyReports[0].Reports[0]
	fmt.Println(
		base.Bold(
			base.ColorizeLine(first.OrderLabel,
				fmt.Sprintf("=== %v %v (%dh %2dm) ===",
					first.OrderLabel, first.SuborderLabel, hours, minutes))))
	//implement list
	for _, dr := range dailyReports {
		printSingleLineForProject(dr)
	}
}

func printSingleLineForProject(dailyReport base.DailyReport) {
	report := dailyReport.Reports[0]
	var training string
	if report.Training {
		training = "Fobi"
	} else {
		training = "    "
	}
	fmt.Println(
		fmt.Sprintf("%v %v  %dh %2dm  %v \t%s",
			formatWeekdayGermanShort(dailyReport.Date),
			formatDateGerman(dailyReport.Date),
			report.Hours,
			report.Minutes,
			training,
			report.Comment))
}

func formatDateGerman(date time.Time) string {
	return date.Format("02.01.2006")
}

func formatWeekdayGermanShort(date time.Time) string {
	switch date.Weekday() {
	case time.Monday:
		return "Mo"
	case time.Tuesday:
		return "Di"
	case time.Wednesday:
		return "Mi"
	case time.Thursday:
		return "Do"
	case time.Friday:
		return "Fr"
	case time.Saturday:
		return "Sa"
	case time.Sunday:
		return "So"
	default:
		return "Fehlerhafter Wochentag"
	}
}

func printSingleLine(index int, report base.DailyReportData) {
	var training string
	if report.Training {
		training = "Fobi"
	} else {
		training = "    "
	}
	fmt.Println(
		base.ColorizeLine(report.OrderLabel,
			fmt.Sprintf("%x. %dh %2dm %v \t%-20s  %-25s\t%s",
				index+1,
				report.Hours,
				report.Minutes,
				training,
				report.OrderLabel,
				report.SuborderLabel,
				report.Comment)))

	/*
		fmt.Printf("%x. %dh %2dm %v \t%s\t%s\t%s\n",
			index+1,
			report.Hours,
			report.Minutes,
			training,
			base.Colorize(report.OrderLabel, report.OrderLabel, 20),
			base.Colorize(report.OrderLabel, report.SuborderLabel, 25),
			report.Comment)
	*/
}

func split(orderLists []base.DailyReport) []base.DailyReport {
	count := 0
	for _, d := range orderLists {
		count = count + len(d.Reports)
	}
	result := make([]base.DailyReport, 0, count)

	for _, d := range orderLists {
		if len(d.Reports) == 1 {
			result = append(result, d)
		} else {
			for _, r := range d.Reports {
				reportList := make([]base.DailyReportData, 0, 1)
				reportList = append(reportList, r)
				dr := base.DailyReport{Date: d.Date,
					Reports: reportList}
				result = append(result, dr)
			}

		}
	}
	return result
}

func compare(a base.DailyReport, b base.DailyReport) bool {
	if a.Reports[0].OrderLabel != b.Reports[0].OrderLabel {
		return a.Reports[0].OrderLabel < b.Reports[0].OrderLabel
	}
	return a.Reports[0].SuborderLabel < b.Reports[0].SuborderLabel
}
