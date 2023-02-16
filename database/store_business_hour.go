package database

import "time"

type StoreBusinessHour struct {
	ID             int64     `pg:"id"`
	StoreID        int64     `pg:"store_id"`
	DayOfWeek      int64     `pg:"day_of_week"`
	StartTimeLocal time.Time `pg:"start_time_local"`
	EndTimeLocal   time.Time `pg:"end_time_local"`
}
