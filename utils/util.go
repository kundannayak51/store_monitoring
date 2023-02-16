package utils

import (
	"time"
)

const CurrentTime = "2023-01-24 09:07:26.441407 UTC"

func ConvertUTCToLocal(utcTimestamp, timezoneStr string) (time.Time, string, error) {
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

func GetTimeOfXDaysBefore(currentTime string, days int) (string, error) {
	givenTime, err := time.Parse(time.RFC3339Nano, currentTime)
	if err != nil {
		return "", err
	}

	oneWeekAgo := givenTime.AddDate(0, 0, (days * (-1)))

	return oneWeekAgo.Format(time.RFC3339Nano), nil
}
