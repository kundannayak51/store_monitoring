package utils

import (
	"github.com/store_monitoring/database"
	"github.com/store_monitoring/entities"
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
