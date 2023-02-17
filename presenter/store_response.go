package presenter

import "github.com/gin-gonic/gin"

type GetReportResponse struct {
}

type TriggerReportData struct {
	ReportId string
}

type GeneralErrorResponse struct {
	Message string `json:"message"`
}

func GetTriggerReportResponse(reportId string) TriggerReportData {
	return TriggerReportData{
		ReportId: reportId,
	}
}

func HandleGeneralErrorResponse(ctx *gin.Context, err error) *GeneralErrorResponse {
	return &GeneralErrorResponse{
		Message: err.Error(),
	}
}
