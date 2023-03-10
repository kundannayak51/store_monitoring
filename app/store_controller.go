package app

import (
	"bytes"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"github.com/store_monitoring/presenter"
	"github.com/store_monitoring/services"
	"github.com/store_monitoring/utils"
	"net/http"
	"strconv"
)

type StoreController struct {
	StoreService services.StoreServiceInterface
}

func NewStoreController(storeService services.StoreServiceInterface) *StoreController {
	return &StoreController{
		StoreService: storeService,
	}
}

func (con *StoreController) TriggerReport(c *gin.Context) {
	ctx := utils.GetValueOnlyRequestContext(c)
	reportId, err := con.StoreService.TriggerReportGeneration(ctx)
	if err != nil {
		presenter.HandleGeneralErrorResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, presenter.GetTriggerReportResponse(reportId))
}

func (con *StoreController) GetCSVReport(c *gin.Context) {
	reportId := c.Param("report_id")
	if len(reportId) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Report Id",
		})
		return
	}
	ctx := utils.GetValueOnlyRequestContext(c)
	reports, err := con.StoreService.GetCSVData(ctx, reportId)
	if err != nil {
		presenter.HandleGeneralErrorResponse(c, err)
	}
	if len(reports) == 0 {
		c.String(http.StatusOK, "Running")
		return
	}
	buffer := &bytes.Buffer{}

	writer := csv.NewWriter(buffer)
	writer.Write([]string{"store_id", "uptime_last_hour(%)", "uptime_last_day(%)", "uptime_last_week(%)", "downtime_last_hour(%)", "downtime_last_day(%)", "downtime_last_week(%)"})

	for _, r := range reports {
		writer.Write([]string{strconv.FormatInt(r.StoreId, 10), utils.ConvertFloat64ToString(r.UptimeLastHour), utils.ConvertFloat64ToString(r.UptimeLastDay), utils.ConvertFloat64ToString(r.UptimeLastWeek), utils.ConvertFloat64ToString(r.DowntimeLastHour), utils.ConvertFloat64ToString(r.DowntimeLastDay), utils.ConvertFloat64ToString(r.DowntimeLastWeek)})
	}

	writer.Flush()

	// Set the response header to indicate that this is a CSV file
	c.Header("Content-Disposition", "attachment; filename=report.csv")
	c.Data(http.StatusOK, "text/csv", buffer.Bytes())
	c.String(http.StatusOK, "Completed")

}
