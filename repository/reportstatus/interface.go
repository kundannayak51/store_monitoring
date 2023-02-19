package reportstatus

import (
	"context"
	"github.com/store_monitoring/entities"
)

type ReportStatusRepo interface {
	InsertReportStatus(ctx context.Context, reportID string, status string) error
	UpdateStatusForReportId(ctx context.Context, reportId string, status string) (int64, error)
	GetReportStatus(ctx context.Context, reportId string) (*entities.ReportStatus, error)
}
