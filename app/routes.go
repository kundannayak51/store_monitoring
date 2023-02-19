package app

import (
	"database/sql"
	"github.com/store_monitoring/repository/report"
	"github.com/store_monitoring/repository/reportstatus"
	"github.com/store_monitoring/repository/storebusinesshour"
	"github.com/store_monitoring/repository/storestatus"
	"github.com/store_monitoring/repository/storetimezone"
	"github.com/store_monitoring/services"
)

func SetupRoutes(db *sql.DB) {
	storeBusinessHourRepo := storebusinesshour.NewStoreBusinessHourRepository(db)
	storeStatusRepo := storestatus.NewStoreStatusRepository(db)
	storeTimezoneRepo := storetimezone.NewStoreTimezoneRepository(db)
	reportStatusRepo := reportstatus.NewReportStatusRepository(db)
	reportRepo := report.NewReportRepository(db)

	storeService := services.NewService(storeBusinessHourRepo, storeStatusRepo, storeTimezoneRepo, reportStatusRepo, reportRepo)

	storeController := NewStoreController(*storeService)

	Router.POST("/trigger_report", storeController.TriggerReport)
	Router.GET("/get_report/:report_id", storeController.GetCSVReport)
}
