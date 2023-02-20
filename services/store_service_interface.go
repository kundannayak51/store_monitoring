package services

import (
	"context"
	"github.com/store_monitoring/entities"
)

type StoreServiceInterface interface {
	GetCSVData(ctx context.Context, reportId string) ([]entities.Report, error)
	TriggerReportGeneration(ctx context.Context) (string, error)
}
