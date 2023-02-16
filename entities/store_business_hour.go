package entities

import "time"

type StoreBusinessHour struct {
	StoreID        int64
	DayOfWeek      int64
	StartLocalTime time.Time
	EndLocalTime   time.Time
}
