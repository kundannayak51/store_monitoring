package entities

import "time"

type BusinessHour struct {
	StoreID        int64
	DayOfWeek      int
	StartLocalTime time.Time
	EndLocalTime   time.Time
}
