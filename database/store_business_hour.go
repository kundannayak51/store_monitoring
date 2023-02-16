package database

type StoreBusinessHour struct {
	ID          int64
	StoreID     int64
	DayOfWeek   string
	OpeningTime string
	ClosingTime string
}
