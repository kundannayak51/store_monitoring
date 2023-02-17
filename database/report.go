package database

type Report struct {
	ReportId         string  `pg:"report_id"`
	StoreId          int64   `pg:"store_id"`
	UptimeLastHour   float64 `pg:"uptime_last_hour"`
	UptimeLastDay    float64 `pg:"uptime_last_day"`
	UptimeLastWeek   float64 `pg:"uptime_last_week"`
	DowntimeLastHour float64 `pg:"downtime_last_hour"`
	DowntimeLastDay  float64 `pg:"downtime_last_day"`
	DowntimeLastWeek float64 `pg:"downtime_last_week"`
}
