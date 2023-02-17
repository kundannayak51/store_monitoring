package services

import (
	"context"
	"fmt"
	"github.com/store_monitoring/entities"
	"github.com/store_monitoring/repository/report"
	"github.com/store_monitoring/repository/reportstatus"
	"github.com/store_monitoring/repository/storebusinesshour"
	"github.com/store_monitoring/repository/storestatus"
	"github.com/store_monitoring/repository/storetimezone"
	"github.com/store_monitoring/utils"
	"time"
)

type StoreService struct {
	storeBusinessHourRepo storebusinesshour.StoreBusinessHourRepository
	storeStatusRepo       storestatus.StoreStatusRepository
	storeTimezoneRepo     storetimezone.StoreTimezoneRepository
	reportStatusRepo      reportstatus.ReportStatusRepository
	reportRepo            report.ReportRepository
}

func NewService(ctx context.Context, storeBusinessHourRepo storebusinesshour.StoreBusinessHourRepository, storeStatusRepo storestatus.StoreStatusRepository, storeTimezoneRepo storetimezone.StoreTimezoneRepository, reportStatusRepo reportstatus.ReportStatusRepository, reportRepo report.ReportRepository) *StoreService {
	return &StoreService{
		storeBusinessHourRepo: storeBusinessHourRepo,
		storeStatusRepo:       storeStatusRepo,
		storeTimezoneRepo:     storeTimezoneRepo,
		reportStatusRepo:      reportStatusRepo,
		reportRepo:            reportRepo,
	}
}

func (s *StoreService) GetCSVData(ctx context.Context, reportId string) ([]entities.Report, error) {
	// Query the database for report status
	//var reportStatus database.ReportStatus
	//err := s.db.QueryRow("SELECT * FROM ReportStatus WHERE report_id = $1", reportID).Scan(&reportStatus.ReportId, &reportStatus.Status)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Report ID not found"})
	//	return
	//}
	reportStatus, err := s.reportStatusRepo.GetReportStatus(ctx, reportId)
	if err != nil {
		return nil, err
	}

	// Check if report is completed
	if reportStatus.Status != "Completed" {
		return []entities.Report{}, nil
	}

	// Query the database for report data
	reports, err := s.reportRepo.GetReportsForReportId(ctx, reportId)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

func (s *StoreService) TriggerReportGeneration(ctx context.Context) (string, error) {
	storeIds, err := s.storeTimezoneRepo.GetAllStores()
	if err != nil {
		return "", err
	}
	reportId := utils.GenerateReportId()
	err = s.reportStatusRepo.InsertReportStatus(ctx, reportId, "Running")
	if err != nil {
		return "", err
	}
	go s.triggerReportGenerationForEachStore(ctx, storeIds, reportId)
	return reportId, nil
}

func (s *StoreService) triggerReportGenerationForEachStore(ctx context.Context, storeIds []int64, reportId string) {
	for _, id := range storeIds {
		err := s.GenerateAndStoreReportForStoreId(ctx, id, reportId)
		if err != nil {
			//log error
		}
	}
	rowsAffected, err := s.reportStatusRepo.UpdateStatusForReportId(ctx, reportId, "Completed")
	if err != nil {
		//log error
	}
	//log rows affected
	//TODO: remove this with log
	fmt.Sprintf("Number of rows updated:", rowsAffected)
}

func (s *StoreService) GenerateAndStoreReportForStoreId(ctx context.Context, storeId int64, reportId string) error {
	endTime := utils.CurrentTime
	startTime, err := utils.GetTimeOfXDaysBefore(endTime, 7)

	timeZone, err := s.storeTimezoneRepo.GetTimezoneForStore(storeId)
	if err != nil {
		return err
	}

	localStartTime, _, err := utils.ConvertUTCStrToLocal(startTime, timeZone)
	if err != nil {
		return err
	}
	localEndTime, _, err := utils.ConvertUTCStrToLocal(endTime, timeZone)
	if err != nil {
		return err
	}
	//Fetching Business Hours of a Store within a time range of one week
	storeBusinessHours, err := s.storeBusinessHourRepo.GetBusinessHoursInTimeRange(ctx, storeId, localStartTime, localEndTime)
	if err != nil {
		return err
	}
	//If a week day is missing in storeBusinessHours, enriching that day with startTime "00:00:00" and endTime "23:59:59"
	storeDayTimeMapping, err := s.enrichBusinessHoursAndReturnDayTimeMapping(ctx, storeId, &storeBusinessHours)
	if err != nil {
		return err
	}

	storeStatuses, err := s.storeStatusRepo.GetStoreStatusInTimeRange(ctx, storeId, startTime, endTime)
	if err != nil {
		return err
	}

	report, err := s.calculateWeeklyObservationAndGererateReport(ctx, storeId, &storeStatuses, &storeBusinessHours, timeZone, storeDayTimeMapping, reportId)
	if err != nil {
		return err
	}

	err = s.reportRepo.InsertReport(ctx, report)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreService) calculateWeeklyObservationAndGererateReport(ctx context.Context, storeId int64, storeStatuses *[]entities.StoreStatus, storeBusinessHours *[]entities.StoreBusinessHour, timeZone string, storeDayTimeMapping map[int64]entities.StartEndTime, reportId string) (*entities.Report, error) {
	statusMap := make(map[int64]map[int64]string, 0)

	for _, businessHour := range *storeBusinessHours {
		startTime := businessHour.StartLocalTime
		endTime := businessHour.EndLocalTime

		day := businessHour.DayOfWeek

		totalHourChunks := s.calculateTotalChunks(startTime, endTime)
		for i := 0; i < totalHourChunks; i++ {
			statusMap[day][int64(i)] = "None"
		}
	}

	for _, storeStatus := range *storeStatuses {
		localDateTime, dayStr, err := utils.ConvertUTCToLocal(storeStatus.Timestamp, timeZone)
		if err != nil {
			return nil, err
		}
		day := utils.GetDayMapping(dayStr)
		localTime := utils.GetTimeFromTime(localDateTime)

		chunkNumber := s.getChunkNumberFromEnd(storeDayTimeMapping[day].StartTime, storeDayTimeMapping[day].EndTime, localTime)
		statusMap[day][chunkNumber] = storeStatus.Status
	}

	statusMap = s.enrichStatusMapWithNearestStatus(ctx, &statusMap)

	currentTime := utils.CurrentTime
	_, dayStr, err := utils.ConvertUTCStrToLocal(currentTime, timeZone)
	if err != nil {
		return nil, err
	}
	day := utils.GetDayMapping(dayStr)

	weekelyObservation := s.createWeeklyObservation(ctx, storeId, &statusMap, day)

	report := s.generateWeeklyReport(ctx, storeId, weekelyObservation, day, reportId)
	return report, nil
}

func (s *StoreService) generateWeeklyReport(ctx context.Context, storeId int64, observations *entities.Observation, currentDay int64, reportId string) *entities.Report {
	var uptimeLastHour, uptimeLastDay, uptimeLastWeek, downtimeLastHour, downtimeLastDay, downtimeLastWeek float64

	if observations.IsLastHourActive {
		uptimeLastHour = 100.0
		downtimeLastHour = 0.0
	} else {
		downtimeLastHour = 100.0
		uptimeLastHour = 0.0
	}

	totalWeeklyChunks, totalLastDayChunks, totalActiveWeklyChunks, totalActiveLastDayChunks := 0, 0, 0, 0

	for _, observation := range observations.WeekReport {
		totalWeeklyChunks += int(observation.TotalChunks)
		totalActiveWeklyChunks += int(observation.ActiveChunks)
		if currentDay == observation.Day {
			totalLastDayChunks = int(observation.TotalChunks)
			totalActiveLastDayChunks = int(observation.ActiveChunks)
		}
	}

	uptimeLastDay = float64(totalActiveLastDayChunks) / float64(totalLastDayChunks) * 100
	downtimeLastDay = float64(totalLastDayChunks-totalActiveLastDayChunks) / float64(totalLastDayChunks) * 100

	uptimeLastWeek = float64(totalActiveWeklyChunks) / float64(totalWeeklyChunks) * 100
	downtimeLastWeek = float64(totalWeeklyChunks-totalActiveWeklyChunks) / float64(totalWeeklyChunks) * 100

	return &entities.Report{
		ReportId:         reportId,
		StoreId:          storeId,
		UptimeLastDay:    uptimeLastDay,
		UptimeLastHour:   uptimeLastHour,
		UptimeLastWeek:   uptimeLastWeek,
		DowntimeLastDay:  downtimeLastDay,
		DowntimeLastHour: downtimeLastHour,
		DowntimeLastWeek: downtimeLastWeek,
	}
}

func (s *StoreService) createWeeklyObservation(ctx context.Context, storeId int64, statusMap *map[int64]map[int64]string, currentDay int64) *entities.Observation {
	weeklyObservation := make([]entities.WeeklyObservation, 0)
	isLastHourActive := false

	for key, val := range *statusMap {
		totalChunks := len(val)
		activeStatus := 0

		for chunk, status := range val {
			if status == "active" {
				activeStatus++
				if key == currentDay && chunk == 0 {
					isLastHourActive = true
				}
			}
		}
		dayObservation := entities.WeeklyObservation{
			Day:          key,
			TotalChunks:  int64(totalChunks),
			ActiveChunks: int64(activeStatus),
		}
		weeklyObservation = append(weeklyObservation, dayObservation)
	}
	return &entities.Observation{
		StoreId:          storeId,
		WeekReport:       weeklyObservation,
		IsLastHourActive: isLastHourActive,
	}
}

func (s *StoreService) enrichStatusMapWithNearestStatus(ctx context.Context, statusMap *map[int64]map[int64]string) map[int64]map[int64]string {
	resultMap := make(map[int64]map[int64]string, 0)
	for key, val := range *statusMap {
		lastStatus := "inactive"
		lastStatusIndex := 0
		countActive := 0
		countInactive := 0
		tempMap := make(map[int64]string, 0)
		for i := 0; i < len(val); i++ {
			if val[int64(i)] == "inactive" {
				countInactive++
			} else if val[int64(i)] == "active" {
				countActive++
			}
		}

		for i := 0; i < len(val); i++ {
			if val[int64(i)] != "None" {
				lastStatus = val[int64(i)]
				lastStatusIndex = i
				tempMap[int64(i)] = val[int64(i)]
			} else {
				j := 0
				for j = i + 1; j < len(val); j++ {
					if j-i > i-lastStatusIndex {
						tempMap[int64(i)] = lastStatus
						break
					} else {
						if val[int64(j)] != "None" {
							if j-i == i-lastStatusIndex {
								if countInactive > countActive {
									tempMap[int64(i)] = "inactive"
								} else {
									tempMap[int64(i)] = "active"
								}
								break
							}
							tempMap[int64(i)] = val[int64(j)]
							break
						}
					}
				}
				if j == len(val) {
					if countInactive > countActive {
						tempMap[int64(i)] = "inactive"
					} else {
						tempMap[int64(i)] = "active"
					}
				}
			}
		}
		resultMap[key] = tempMap
	}
	return resultMap
}

func (s *StoreService) enrichBusinessHoursAndReturnDayTimeMapping(ctx context.Context, storeId int64, businessHours *[]entities.StoreBusinessHour) (map[int64]entities.StartEndTime, error) {
	storeDayTimeMapping := make(map[int64]entities.StartEndTime)

	startTime, err := time.Parse("15:04:05", utils.LocalTimeStart)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse("15:04:05", utils.LocalTimeEnd)
	if err != nil {
		return nil, err
	}
	daysPresent := make(map[int64]bool, 0)
	for _, businessHour := range *businessHours {
		dayOfWeek := businessHour.DayOfWeek
		daysPresent[dayOfWeek] = true
		storeDayTimeMapping[dayOfWeek] = entities.StartEndTime{
			StartTime: businessHour.StartLocalTime,
			EndTime:   businessHour.EndLocalTime,
		}
	}

	for i := 0; i <= 6; i++ {
		_, present := daysPresent[int64(i)]
		if !present {
			newBusinessHour := entities.StoreBusinessHour{
				StoreID:        storeId,
				DayOfWeek:      int64(i),
				StartLocalTime: startTime,
				EndLocalTime:   endTime,
			}
			storeDayTimeMapping[int64(i)] = entities.StartEndTime{
				StartTime: startTime,
				EndTime:   endTime,
			}
			*businessHours = append(*businessHours, newBusinessHour)
		}
	}
	return storeDayTimeMapping, nil
}

func (s *StoreService) calculateTotalChunks(startTime, endTime time.Time) int {
	/*start, err := time.Parse("15:04:05", startTime)
	if err != nil {
		log.Fatal(err)
	}

	end, err := time.Parse("15:04:05", endTime)
	if err != nil {
		log.Fatal(err)
	}*/
	duration := endTime.Sub(startTime)
	totalMinutes := int(duration.Minutes())

	if totalMinutes%60 == 0 {
		return totalMinutes / 60
	}
	return totalMinutes/60 + 1
}

func (s *StoreService) getChunkNumberFromEnd(startTime, endTime, inputTime time.Time) int64 {
	// Parse start and end times
	//start, _ := time.Parse("15:04:05", startTime)
	//end, _ := time.Parse("15:04:05", endTime)
	//inputTime, _ := time.Parse("15:04:05", timeStr)

	// Calculate the number of minutes between start time and input time
	minutes := int(endTime.Sub(inputTime).Minutes())

	// Calculate the chunk number
	chunk := minutes / 60

	return int64(chunk)
}
