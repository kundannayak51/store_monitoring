package storestatus

import "database/sql"

type StoreStatusRepository struct {
	db *sql.DB
}

func NewStoreStatusRepository(db *sql.DB) *StoreStatusRepository {
	return &StoreStatusRepository{db}
}
