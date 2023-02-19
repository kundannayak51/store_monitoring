package report

import (
	"context"
	"github.com/store_monitoring/entities"
)

type ReportRepo interface {
	InsertReport(ctx context.Context, report *entities.Report) error
	GetReportsForReportId(ctx context.Context, reportId string) ([]entities.Report, error)
}
