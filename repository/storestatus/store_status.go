package storestatus

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/store_monitoring/database"
	"github.com/store_monitoring/entities"
	"github.com/store_monitoring/utils"
	"time"
)

type StoreStatusRepository struct {
	db *sql.DB
}

func NewStoreStatusRepository(db *sql.DB) *StoreStatusRepository {
	return &StoreStatusRepository{db}
}

const getStoreStatusQuery = `
        SELECT id, store_id, timestamp_utc, status
        FROM store_status
        WHERE store_id = $1
            AND timestamp_utc >= $2
            AND timestamp_utc < $3
    `

func (s *StoreStatusRepository) GetStoreStatusInTimeRange(ctx context.Context, storeId int64, startTimeStr, endTimeStr string) ([]entities.StoreStatus, error) {

	// execute the query and retrieve the rows
	startTime, err := time.Parse("2006-01-02 15:04:05.999999999 MST", startTimeStr)
	if err != nil {
		return nil, err
	}
	endTime, err := time.Parse("2006-01-02 15:04:05.999999999 MST", endTimeStr)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(getStoreStatusQuery, storeId, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	storeStatuses := make([]entities.StoreStatus, 0)
	for rows.Next() {
		var s database.StoreStatus
		err := rows.Scan(&s.ID, &s.StoreID, &s.Timestamp, &s.Status)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		storeStatuses = append(storeStatuses, *utils.ConvertStoreStatusDaoToEntity(&s))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error retrieving rows: %v", err)
	}

	return storeStatuses, nil
}
