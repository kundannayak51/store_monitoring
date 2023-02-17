package storetimezone

import (
	"database/sql"
	"fmt"
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

func (sr *StoreTimezoneRepository) GetAllStores() ([]int64, error) {
	rows, err := sr.db.Query("SELECT store_id FROM timezones")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %v", err)
	}
	defer rows.Close()

	storeIds := make([]int64, 0)
	for rows.Next() {
		var id int64
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		storeIds = append(storeIds, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error retrieving rows: %v", err)
	}

	return storeIds, nil
}
