package services

import (
	"context"
	"github.com/store_monitoring/constants"
	"github.com/store_monitoring/entities"
	"github.com/store_monitoring/repository/report"
	"github.com/store_monitoring/repository/reportstatus"
	"github.com/store_monitoring/repository/storebusinesshour"
	"github.com/store_monitoring/repository/storestatus"
	"github.com/store_monitoring/repository/storetimezone"
	"github.com/store_monitoring/utils"
	"math/rand"
	"time"
)

type StoreService struct {
	storeBusinessHourRepo storebusinesshour.StoreBusinessHourRepo
	storeStatusRepo       storestatus.StoreStatusRepo
	storeTimezoneRepo     storetimezone.StoreTimezoneRepo
	reportStatusRepo      reportstatus.ReportStatusRepo
	reportRepo            report.ReportRepo
}

func NewService(storeBusinessHourRepo storebusinesshour.StoreBusinessHourRepo, storeStatusRepo storestatus.StoreStatusRepo, storeTimezoneRepo storetimezone.StoreTimezoneRepo, reportStatusRepo reportstatus.ReportStatusRepo, reportRepo report.ReportRepo) *StoreService {
	return &StoreService{
		storeBusinessHourRepo: storeBusinessHourRepo,
		storeStatusRepo:       storeStatusRepo,
		storeTimezoneRepo:     storeTimezoneRepo,
		reportStatusRepo:      reportStatusRepo,
		reportRepo:            reportRepo,
	}
}

func (s *StoreService) GetCSVData(ctx context.Context, reportId string) ([]entities.Report, error) {

	reportStatus, err := s.reportStatusRepo.GetReportStatus(ctx, reportId)
	if err != nil {
		return nil, err
	}

	// Check if report is completed
	if reportStatus.Status != constants.STATUS_COMPLETE {
		return []entities.Report{}, nil
	}

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
	err = s.reportStatusRepo.InsertReportStatus(ctx, reportId, constants.STATUS_RUNNING)
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
			return
		}
	}
	_, err := s.reportStatusRepo.UpdateStatusForReportId(ctx, reportId, constants.STATUS_COMPLETE)
	if err != nil {
		//log error
	}
	//log rows affected
	//TODO: remove this with log
}

func (s *StoreService) GenerateAndStoreReportForStoreId(ctx context.Context, storeId int64, reportId string) error {
	endTime := utils.CurrentTime
	startTime, err := utils.GetTimeOfXDaysBefore(endTime, 7)

	//TODO: no need of this, we can fetch timezones along with storeID only
	timeZone, err := s.storeTimezoneRepo.GetTimezoneForStore(storeId)
	if err != nil {
		return err
	}

	//Fetching Business Hours of a Store within a time range of one week
	storeBusinessHours, err := s.storeBusinessHourRepo.GetBusinessHoursInTimeRange(ctx, storeId)
	if err != nil {
		//log error
		return err
	}
	//If a week day is missing in storeBusinessHours, enriching that day with startTime "00:00:00" and endTime "23:59:59"
	storeDayTimeMapping, err := s.enrichBusinessHoursAndReturnDayTimeMapping(ctx, storeId, &storeBusinessHours, timeZone)
	if err != nil {
		//log error
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
		statusMap[day] = make(map[int64]string)
		for i := 0; i < totalHourChunks; i++ {
			statusMap[day][int64(i)] = constants.STATUS_NONE
		}
	}

	for _, storeStatus := range *storeStatuses {
		localTime, dayStr, err := utils.ConvertUTCToLocal(storeStatus.Timestamp, timeZone)
		if err != nil {
			return nil, err
		}
		day := utils.GetDayMapping(dayStr)

		liesBetween, err := utils.CheckUTCTimeLiesBetweenTwoLocalTime(storeDayTimeMapping[day].StartTime, storeDayTimeMapping[day].EndTime, localTime, timeZone)
		if err != nil {
			return nil, err
		}
		if !liesBetween {
			continue
		}

		chunkNumber := s.getChunkNumberFromEnd(storeDayTimeMapping[day].EndTime, localTime)
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
			if status == constants.STATUS_ACTIVE {
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
		lastStatus := constants.STATUS_INACTIVE
		rand.Seed(time.Now().UnixNano())
		rVal := rand.Float64()
		if rVal >= 0.5 {
			lastStatus = constants.STATUS_ACTIVE
		}
		lastStatusIndex := 0
		countActive := 0
		countInactive := 0
		tempMap := make(map[int64]string, 0)
		for i := 0; i < len(val); i++ {
			if val[int64(i)] == constants.STATUS_INACTIVE {
				countInactive++
			} else if val[int64(i)] == constants.STATUS_ACTIVE {
				countActive++
			}
		}

		if countActive == 0 && countInactive == 0 {
			for i := 0; i < len(val); i++ {

				randVal := rand.Float64()
				if randVal >= 0.5 {
					tempMap[int64(i)] = constants.STATUS_ACTIVE
				} else {
					tempMap[int64(i)] = constants.STATUS_INACTIVE
				}
			}

		}

		for i := 0; i < len(val); i++ {
			if val[int64(i)] != constants.STATUS_NONE {
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
						if val[int64(j)] != constants.STATUS_NONE {
							if j-i == i-lastStatusIndex {
								if countInactive > countActive {
									tempMap[int64(i)] = constants.STATUS_INACTIVE
								} else {
									tempMap[int64(i)] = constants.STATUS_ACTIVE
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
						tempMap[int64(i)] = constants.STATUS_INACTIVE
					} else {
						tempMap[int64(i)] = constants.STATUS_ACTIVE
					}
				}
			}
		}
		resultMap[key] = tempMap
	}
	return resultMap
}

func (s *StoreService) enrichBusinessHoursAndReturnDayTimeMapping(ctx context.Context, storeId int64, businessHours *[]entities.StoreBusinessHour, timezone string) (map[int64]entities.StartEndTime, error) {
	storeDayTimeMapping := make(map[int64]entities.StartEndTime)

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
				StartLocalTime: utils.LocalTimeStart,
				EndLocalTime:   utils.LocalTimeEnd,
			}
			storeDayTimeMapping[int64(i)] = entities.StartEndTime{
				StartTime: utils.LocalTimeStart,
				EndTime:   utils.LocalTimeEnd,
			}
			*businessHours = append(*businessHours, newBusinessHour)
		}
	}
	return storeDayTimeMapping, nil
}

func (s *StoreService) calculateTotalChunks(startTimeStr, endTimeStr string) int {
	startTime, err := time.Parse("15:04:05", startTimeStr)
	if err != nil {
		return -1
	}

	endTime, err := time.Parse("15:04:05", endTimeStr)
	if err != nil {
		return -1
	}
	duration := endTime.Sub(startTime)
	totalMinutes := int(duration.Minutes())

	if totalMinutes%60 == 0 {
		return totalMinutes / 60
	}
	return totalMinutes/60 + 1
}

func (s *StoreService) getChunkNumberFromEnd(endTimeStr, inputTimeStr string) int64 {

	inputTime, err := time.Parse("15:04:05", inputTimeStr)
	if err != nil {
		return -1
	}

	endTime, err := time.Parse("15:04:05", endTimeStr)
	if err != nil {
		return -1
	}

	minutes := int(endTime.Sub(inputTime).Minutes())

	chunk := minutes / 60

	return int64(chunk)
}
