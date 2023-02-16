package store_monitoring

import (
	"fmt"
	"github.com/store_monitoring/database"
	"github.com/store_monitoring/repository/storebusinesshour"
)

func main() {
	// initialize the database connection
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// create a new instance of the store business hour repository
	storeBusinessHourRepo := storebusinesshour.NewStoreBusinessHourRepository(db)

	// call the GetStoreBusinessHour method to get the store business hours
	storeBusinessHour, err := storeBusinessHourRepo.GetStoreBusinessHour(storeID)
	if err != nil {
		fmt.Println("Error getting store business hours:", err)
		return
	}

	// print the store business hours
	fmt.Println(storeBusinessHour)
}
