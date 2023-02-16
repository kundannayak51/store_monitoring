package storetimezone

import (
	"database/sql"
)

type StoreTimezoneRepository struct {
	db *sql.DB
}

func NewStoreTimezoneRepository(db *sql.DB) *StoreTimezoneRepository {
	return &StoreTimezoneRepository{db}
}

func (sr *StoreTimezoneRepository) GetTimezoneForStore(storeId int64) (string, error) {
	var timezone string
	err := sr.db.QueryRow("SELECT timezone_str FROM timezones WHERE store_id = $1", storeId).Scan(&timezone)
	if err != nil {
		return "", err
	}
	return timezone, nil
}
