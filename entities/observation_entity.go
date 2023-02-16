package entities

type Observation struct {
	StoreId          int64
	WeekReport       []WeeklyObservation
	IsLastHourActive bool
}

type WeeklyObservation struct {
	Day          int64
	TotalChunks  int64
	ActiveChunks int64
}
