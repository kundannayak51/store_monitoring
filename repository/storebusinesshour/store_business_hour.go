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

func (s *StoreBusinessHourRepository) GetStoreBusinessHour(storeID int64) (*database.StoreBusinessHour, error) {
	query := `SELECT * FROM store_business_hour WHERE store_id = $storeID`
	row := s.db.QueryRow(query, storeID)

	var storeBusinessHour database.StoreBusinessHour
	err := row.Scan(&storeBusinessHour.ID, &storeBusinessHour.StoreID, &storeBusinessHour.DayOfWeek, &storeBusinessHour.StartTimeLocal, &storeBusinessHour.EndTimeLocal)
	if err != nil {
		return nil, err
	}

	return &storeBusinessHour, nil
}

func (s *StoreBusinessHourRepository) GetBusinessHoursInTimeRange(ctx context.Context, storeID int64, timezone string) ([]entities.StoreBusinessHour, error) {
	// Convert the local start and end times to UTC
	//startTimeUTC := localStartTime.UTC()
	//endTimeUTC := localEndTime.UTC()

	// Query the database for the business hours
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
		businessHours = append(businessHours, *utils.ConvertStoreBusinessHourDaoToEntity(&bh, timezone))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return businessHours, nil
}
