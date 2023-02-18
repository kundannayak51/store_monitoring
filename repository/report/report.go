package report

import (
	"context"
	"database/sql"
	"github.com/store_monitoring/database"
	"github.com/store_monitoring/entities"
	"github.com/store_monitoring/utils"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db}
}

const insertReportQuery = "INSERT INTO report (report_id, store_id, uptime_last_hour, uptime_last_day, uptime_last_week, downtime_last_hour, downtime_last_day, downtime_last_week) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
const getReportQuery = "SELECT report_id, store_id, uptime_last_hour, uptime_last_day, uptime_last_week, downtime_last_hour, downtime_last_day, downtime_last_week FROM report WHERE report_id = $1"

func (r *ReportRepository) InsertReport(ctx context.Context, report *entities.Report) error {
	_, err := r.db.Exec(insertReportQuery, report.ReportId, report.StoreId, report.UptimeLastHour, report.UptimeLastDay, report.UptimeLastWeek, report.DowntimeLastHour, report.DowntimeLastDay, report.DowntimeLastWeek)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReportRepository) GetReportsForReportId(ctx context.Context, reportId string) ([]entities.Report, error) {
	rows, err := r.db.Query(getReportQuery, reportId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reports := make([]entities.Report, 0)
	for rows.Next() {
		var r database.Report
		err := rows.Scan(&r.ReportId, &r.StoreId, &r.UptimeLastHour, &r.UptimeLastDay, &r.UptimeLastWeek, &r.DowntimeLastHour, &r.DowntimeLastDay, &r.DowntimeLastWeek)
		if err != nil {
			return nil, err
		}
		reports = append(reports, *utils.ConvertReportDaoToEntity(&r))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}
