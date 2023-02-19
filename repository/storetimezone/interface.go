package storetimezone

type StoreTimezoneRepo interface {
	GetTimezoneForStore(storeId int64) (string, error)
	GetAllStores() ([]int64, error)
}
