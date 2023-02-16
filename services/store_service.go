package services

import (
	"context"
	"github.com/store_monitoring/repository/storebusinesshour"
	"github.com/store_monitoring/repository/storestatus"
	"github.com/store_monitoring/repository/storetimezone"
	"github.com/store_monitoring/utils"
)

type StoreService struct {
	storeBusinessHourRepo storebusinesshour.StoreBusinessHourRepository
	storeStatusRepo       storestatus.StoreStatusRepository
	storeTimezoneRepo     storetimezone.StoreTimezoneRepository
}

func NewService(ctx context.Context, storeBusinessHourRepo storebusinesshour.StoreBusinessHourRepository, storeStatusRepo storestatus.StoreStatusRepository, storeTimezoneRepo storetimezone.StoreTimezoneRepository) *StoreService {
	return &StoreService{
		storeBusinessHourRepo: storeBusinessHourRepo,
		storeStatusRepo:       storeStatusRepo,
		storeTimezoneRepo:     storeTimezoneRepo,
	}
}

func (s *StoreService) GenerateReportForStoreId(ctx context.Context, storeId int64) {
	endTime := utils.CurrentTime
	startTime, err := utils.GetTimeOfXDaysBefore(endTime, 7)

	timeZone, err := s.storeTimezoneRepo.GetTimezoneForStore(storeId)
	if err != nil {
		return
	}

	localStartTime, startDay, err := utils.ConvertUTCToLocal(startTime, timeZone)
	if err != nil {
		return
	}
	localEndTime, endDay, err := utils.ConvertUTCToLocal(endTime, timeZone)
	if err != nil {
		return
	}
	storeBusinessHours, err := s.storeBusinessHourRepo.GetStoreBusinessHourInTimeRange(storeId)
}