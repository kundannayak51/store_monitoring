package database

import "time"

type StoreStatus struct {
	ID        int64
	StoreID   int64
	Timestamp time.Time
	Status    string
}
