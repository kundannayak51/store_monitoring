package database

type ReportStatus struct {
	ReportId string `pg:"report_id"`
	Status   string `pg:"status"`
}
