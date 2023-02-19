package storebusinesshour

import (
	"context"
	"database/sql"
	"github.com/store_monitoring/database"
	"github.com/store_monitoring/entities"
	"github.com/store_monitoring/utils"
)

type StoreBusinessHourRepository struct {
	db *sql.DB
}

func NewStoreBusinessHourRepository(db *sql.DB) *StoreBusinessHourRepository {
	return &StoreBusinessHourRepository{db}
}

// Define the query to fetch business hours for a store within a time range
const getBusinessHoursQuery = `
	SELECT id, store_id, day_of_week, start_time_local, end_time_local
	FROM store_business_hours
	WHERE store_id = $1
`

func (s *StoreBusinessHourRepository) GetBusinessHoursInTimeRange(ctx context.Context, storeID int64) ([]entities.StoreBusinessHour, error) {
	rows, err := s.db.Query(getBusinessHoursQuery, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and parse the results into a slice of BusinessHour structs
	businessHours := make([]entities.StoreBusinessHour, 0)
	for rows.Next() {
		var bh database.StoreBusinessHour
		err := rows.Scan(&bh.ID, &bh.StoreID, &bh.DayOfWeek, &bh.StartTimeLocal, &bh.EndTimeLocal)
		if err != nil {
			return nil, err
		}
		businessHours = append(businessHours, *utils.ConvertStoreBusinessHourDaoToEntity(&bh))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return businessHours, nil
}
