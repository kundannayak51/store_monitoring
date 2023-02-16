package entities

import "time"

type StoreStatus struct {
	StoreID   int64
	Timestamp time.Time
	Status    string
}
