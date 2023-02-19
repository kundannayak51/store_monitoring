package storestatus

import (
	"context"
	"github.com/store_monitoring/entities"
)

type StoreStatusRepo interface {
	GetStoreStatusInTimeRange(ctx context.Context, storeId int64, startTimeStr, endTimeStr string) ([]entities.StoreStatus, error)
}
