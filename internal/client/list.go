package client

import (
	"encoding/json"
	"fmt"
	"schnegge/internal/base"
	"schnegge/internal/config"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ReadDailyReports(cfg config.Config) []base.DailyReport {
	if !checkConfig(cfg) {
		return nil
	}
	var wg sync.WaitGroup
	date, _ := cfg.GetValue(config.Date)
	dates := CalculateDate(date)
	var dailyReports []base.DailyReport
	for i, date := range dates {
		wg.Add(1)
		dailyReports = append(dailyReports, base.DailyReport{Date: date, Reports: make([]base.DailyReportData, 0)})
		go func(i int, d time.Time) {
			readSingleDay(cfg, d, &dailyReports[i])
			defer wg.Done()
		}(i, date)
	}
	wg.Wait()
	base.Log.Printf("%+v\n", dailyReports)
	return dailyReports
	/*
		// Code for debugging answers
		scanner := bufio.NewScanner(resp.Body)
		for i := 0; scanner.Scan() && i < 50; i++ {
			base.Log.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			base.Log.Panic(err)
		}
		return nil
	*/
}

func readSingleDay(cfg config.Config, date time.Time, into *base.DailyReport) {
	dateString := FormatDateForSalat(date)
	resp, err := doGetRequest(cfg, "/rest/daily-reports/list?refDate="+dateString)
	if err != nil {
		base.Log.Panic(err)
	}
	defer resp.Body.Close()

	base.Log.Println("Response status:", resp.Status)

	err = json.NewDecoder(resp.Body).Decode(&into.Reports)
	if err != nil {
		base.Log.Panic(err)
	}
	base.Log.Printf("%+v\n", *into)
}

func FormatDateForSalat(date time.Time) string {
	return date.Format("2006-01-02")
}

func CalculateDate(dateString string) []time.Time {
	var date time.Time
	var dates []time.Time
	var err error
	switch {
	case strings.HasPrefix(dateString, "-"):
		dates = calculateMinusDays(dateString)
	case strings.EqualFold(dateString, "monat"):
		dates = calculateThisMonth()
	case strings.EqualFold(dateString, "januar") || strings.EqualFold(dateString, "january"):
		dates = calculateMonth(time.January)
	case strings.EqualFold(dateString, "februar") || strings.EqualFold(dateString, "february"):
		dates = calculateMonth(time.February)
	case strings.EqualFold(dateString, "märz") || strings.EqualFold(dateString, "maerz") || strings.EqualFold(dateString, "march"):
		dates = calculateMonth(time.March)
	case strings.EqualFold(dateString, "april"):
		dates = calculateMonth(time.April)
	case strings.EqualFold(dateString, "mai") || strings.EqualFold(dateString, "may"):
		dates = calculateMonth(time.May)
	case strings.EqualFold(dateString, "juni") || strings.EqualFold(dateString, "june"):
		dates = calculateMonth(time.June)
	case strings.EqualFold(dateString, "juli") || strings.EqualFold(dateString, "july"):
		dates = calculateMonth(time.July)
	case strings.EqualFold(dateString, "august"):
		dates = calculateMonth(time.August)
	case strings.EqualFold(dateString, "september"):
		dates = calculateMonth(time.September)
	case strings.EqualFold(dateString, "oktober") || strings.EqualFold(dateString, "october"):
		dates = calculateMonth(time.October)
	case strings.EqualFold(dateString, "november"):
		dates = calculateMonth(time.November)
	case strings.EqualFold(dateString, "dezember") || strings.EqualFold(dateString, "december"):
		dates = calculateMonth(time.December)
	case strings.Contains(dateString, "."):
		date, err = time.Parse("02.01.2006", dateString)
	case strings.Contains(dateString, "-"):
		date, err = time.Parse("2006-01-02", dateString)
	case strings.EqualFold(dateString, "heute") || strings.EqualFold(dateString, "today") || dateString == "":
		date = time.Now()
	case strings.EqualFold(dateString, "gestern") || strings.EqualFold(dateString, "yesterday"):
		date = time.Now().AddDate(0, 0, -1)
	case strings.EqualFold(dateString, "vorgestern"):
		date = time.Now().AddDate(0, 0, -2)
	case strings.EqualFold(dateString, "morgen"):
		date = time.Now().AddDate(0, 0, 1)
	case strings.EqualFold(dateString, "montag") || strings.EqualFold(dateString, "monday"):
		date = calculateWeekday(time.Monday)
	case strings.EqualFold(dateString, "dienstag") || strings.EqualFold(dateString, "tuesday"):
		date = calculateWeekday(time.Tuesday)
	case strings.EqualFold(dateString, "mittwoch") || strings.EqualFold(dateString, "wednesday"):
		date = calculateWeekday(time.Wednesday)
	case strings.EqualFold(dateString, "donnerstag") || strings.EqualFold(dateString, "thursday"):
		date = calculateWeekday(time.Thursday)
	case strings.EqualFold(dateString, "freitag") || strings.EqualFold(dateString, "friday"):
		date = calculateWeekday(time.Friday)
	case strings.EqualFold(dateString, "samstag") || strings.EqualFold(dateString, "saturday"):
		date = calculateWeekday(time.Saturday)
	case strings.EqualFold(dateString, "sonntag") || strings.EqualFold(dateString, "sunday"):
		date = calculateWeekday(time.Sunday)
	case strings.EqualFold(dateString, "woche") || strings.EqualFold(dateString, "week"):
		dates = calculateWeek(0)
	case strings.EqualFold(dateString, "vorwoche"):
		dates = calculateWeek(-7)
	default:
		base.Log.Panic("Error, no valid DatumVar found: " + dateString)
	}
	if err != nil {
		fmt.Println("Error during parsing DatumVar: " + dateString)
		base.Log.Panic(err)
	}
	base.Log.Println("DatumVar: ", dateString, " - ", date.Format(time.UnixDate))
	if len(dates) == 0 {
		dates = append(dates, date)
	}
	return dates
}

func calculateWeekday(weekday time.Weekday) time.Time {
	var date = time.Now()
	var offset = int(weekday) - int(date.Weekday())
	if offset > 0 {
		offset = offset - 7
	}
	date = date.AddDate(0, 0, offset)
	return date
}

func calculateMinusDays(days string) []time.Time {
	var dates []time.Time
	i, err := strconv.Atoi(days)
	if err != nil {
		base.Log.Panic(err)
	}
	for ; i <= 0; i++ {
		dates = append(dates, time.Now().AddDate(0, 0, i))
	}
	return dates
}

func calculateThisMonth() []time.Time {
	minusDays := -1 * (time.Now().Day() - 1)
	return calculateMinusDays(strconv.Itoa(minusDays))
}

func calculateMonth(month time.Month) []time.Time {
	now := time.Now()
	date := time.Date(now.Year(), month, 1, 0, 0, 0, 0, now.Location())
	if date.After(now) {
		date = date.AddDate(-1, 0, 0)
	}

	var dates []time.Time
	for {
		dates = append(dates, date)
		date = date.AddDate(0, 0, 1)
		if date.Day() == 1 {
			break
		}
	}
	return dates
}

func calculateWeek(daysOffset int) []time.Time {
	now := time.Now()
	var offsetToMonday int = int(-1 * ((now.Weekday() + 6) % 7))
	monday := now.AddDate(0, 0, offsetToMonday+daysOffset)

	var dates []time.Time
	for i := 0; i < 7; i++ {
		dates = append(dates, monday.AddDate(0, 0, i))
	}
	return dates
}
