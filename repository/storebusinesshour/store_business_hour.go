package storebusinesshour

import (
	"database/sql"
	"github.com/store_monitoring/database"
)

type StoreBusinessHourRepository struct {
	db *sql.DB
}

func NewStoreBusinessHourRepository(db *sql.DB) *StoreBusinessHourRepository {
	return &StoreBusinessHourRepository{db}
}

func (s *StoreBusinessHourRepository) GetStoreBusinessHour(storeID int64) (*database.StoreBusinessHour, error) {
	query := `SELECT * FROM store_business_hour WHERE store_id = $storeID`
	row := s.db.QueryRow(query, storeID)

	var storeBusinessHour database.StoreBusinessHour
	err := row.Scan(&storeBusinessHour.ID, &storeBusinessHour.StoreID, &storeBusinessHour.DayOfWeek, &storeBusinessHour.OpeningTime, &storeBusinessHour.ClosingTime)
	if err != nil {
		return nil, err
	}

	return &storeBusinessHour, nil
}
