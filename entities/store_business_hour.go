package entities

type StoreBusinessHour struct {
	StoreID        int64
	DayOfWeek      int64
	StartLocalTime string
	EndLocalTime   string
}
