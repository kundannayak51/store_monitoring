package reportstatus

import (
	"context"
	"database/sql"
	"github.com/store_monitoring/database"
	"github.com/store_monitoring/entities"
	"github.com/store_monitoring/utils"
)

type ReportStatusRepository struct {
	db *sql.DB
}

func NewReportStatusRepository(db *sql.DB) *ReportStatusRepository {
	return &ReportStatusRepository{db}
}

func (r *ReportStatusRepository) InsertReportStatus(ctx context.Context, reportID string, status string) error {
	_, err := r.db.Exec("INSERT INTO report_status (report_id, status) VALUES ($1, $2)", reportID, status)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReportStatusRepository) UpdateStatusForReportId(ctx context.Context, reportId string, status string) (int64, error) {
	res, err := r.db.Exec("UPDATE report_status SET status = $1 WHERE report_id = $2", status, reportId)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *ReportStatusRepository) GetReportStatus(ctx context.Context, reportId string) (*entities.ReportStatus, error) {
	var reportStatus database.ReportStatus
	err := r.db.QueryRow("SELECT * FROM report_status WHERE report_id = $1", reportId).Scan(&reportStatus.ReportId, &reportStatus.Status)
	if err != nil {
		return nil, err
	}
	return utils.ConvertReportStatusDaoToEntity(&reportStatus), nil
}
