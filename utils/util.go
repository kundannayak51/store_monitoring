package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/store_monitoring/database"
	"github.com/store_monitoring/entities"
	"strconv"
	"time"
)

const CurrentTime = "2023-01-24 09:07:26.441407 UTC"
const LocalTimeStart = "00:00:00"
const LocalTimeEnd = "23:59:59"

func ConvertUTCStrToLocal(utcTimestamp, timezoneStr string) (time.Time, string, error) {
	utcTime, err := time.Parse("2006-01-02 15:04:05.000000 MST", utcTimestamp)
	if err != nil {
		return time.Time{}, "", err
	}

	location, err := time.LoadLocation(timezoneStr)
	if err != nil {
		return time.Time{}, "", err
	}

	localTime := utcTime.In(location)
	dayOfWeek := localTime.Weekday().String()

	return localTime, dayOfWeek, nil
}

func ConvertUTCToLocal(utcTime time.Time, timezone string) (time.Time, string, error) {
	// Load the specified timezone
	location, err := time.LoadLocation(timezone)
	if err != nil {
		// Handle timezone loading errors
		return time.Time{}, "", err
	}

	// Convert the UTC time to the specified timezone
	localTime := utcTime.In(location)

	// Determine the day of the week
	day := localTime.Weekday().String()

	return localTime, day, nil
}

func GetTimeOfXDaysBefore(currentTime string, days int) (string, error) {
	layout := "2006-01-02 15:04:05.999999999 MST"
	givenTime, err := time.Parse(layout, currentTime)
	if err != nil {
		return "", err
	}

	oneWeekAgo := givenTime.AddDate(0, 0, (days * (-1)))

	return oneWeekAgo.Format(layout), nil
}

func GetTimeFromTime(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
}

func ConvertStoreBusinessHourDaoToEntity(storeBusinessHour *database.StoreBusinessHour) *entities.StoreBusinessHour {
	return &entities.StoreBusinessHour{
		StoreID:        storeBusinessHour.StoreID,
		DayOfWeek:      storeBusinessHour.DayOfWeek,
		StartLocalTime: storeBusinessHour.StartTimeLocal,
		EndLocalTime:   storeBusinessHour.EndTimeLocal,
	}
}

func ConvertStoreStatusDaoToEntity(storeStatus *database.StoreStatus) *entities.StoreStatus {
	return &entities.StoreStatus{
		StoreID:   storeStatus.StoreID,
		Timestamp: storeStatus.Timestamp,
		Status:    storeStatus.Status,
	}
}

func ConvertReportStatusDaoToEntity(reportStatus *database.ReportStatus) *entities.ReportStatus {
	return &entities.ReportStatus{
		ReportId: reportStatus.ReportId,
		Status:   reportStatus.Status,
	}
}

func ConvertReportDaoToEntity(report *database.Report) *entities.Report {
	return &entities.Report{
		ReportId:         report.ReportId,
		StoreId:          report.StoreId,
		UptimeLastHour:   report.UptimeLastHour,
		UptimeLastDay:    report.UptimeLastDay,
		UptimeLastWeek:   report.UptimeLastWeek,
		DowntimeLastHour: report.DowntimeLastHour,
		DowntimeLastDay:  report.DowntimeLastDay,
		DowntimeLastWeek: report.DowntimeLastWeek,
	}
}

func GetDayMapping(day string) int64 {
	switch day {
	case "Monday":
		return 0
	case "Tuesday":
		return 1
	case "Wednesday":
		return 2
	case "Thursday":
		return 3
	case "Friday":
		return 4
	case "Saturday":
		return 5
	case "Sunday":
		return 6
	default:
		return -1
	}
}

func GenerateReportId() string {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	reportId := base64.URLEncoding.EncodeToString(randomBytes)

	return reportId
}

func ConvertFloat64ToString(val float64) string {
	s := strconv.FormatFloat(val, 'f', 6, 64)
	return s
}

type ValueOnlyContext struct{ context.Context }

func (ValueOnlyContext) Deadline() (deadline time.Time, ok bool) { return }
func (ValueOnlyContext) Done() <-chan struct{}                   { return nil }
func (ValueOnlyContext) Err() error                              { return nil }
func GetValueOnlyRequestContext(c *gin.Context) ValueOnlyContext {
	return ValueOnlyContext{Context: c.Request.Context()}
}
