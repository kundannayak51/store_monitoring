package database

import "time"

type StoreStatus struct {
	ID        int64     `pg:"id"`
	StoreID   int64     `pg:"store_id"`
	Timestamp time.Time `pg:"timestamp_utc"`
	Status    string    `pg:"status"`
}
