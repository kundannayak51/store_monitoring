package storebusinesshour

import (
	"context"
	"github.com/store_monitoring/entities"
)

type StoreBusinessHourRepo interface {
	GetBusinessHoursInTimeRange(ctx context.Context, storeID int64) ([]entities.StoreBusinessHour, error)
}
