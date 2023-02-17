package entities

type Report struct {
	ReportId         string
	StoreId          int64
	UptimeLastHour   float64
	UptimeLastDay    float64
	UptimeLastWeek   float64
	DowntimeLastHour float64
	DowntimeLastDay  float64
	DowntimeLastWeek float64
}
