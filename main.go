package main

import (
	"context"
	"github.com/store_monitoring/app"
)

func main() {
	// initialize the database connection
	/*db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// create a new instance of the store business hour repository
	storeId := 8605792781614382846
	storeBusinessHourRepo := storebusinesshour.NewStoreBusinessHourRepository(db)
	storeStatusRepo := storestatus.NewStoreStatusRepository(db)

	// call the GetStoreBusinessHour method to get the store business hours
	//storeBusinessHour, err := storeBusinessHourRepo.GetStoreBusinessHour(int64(storeId))
	//if err != nil {
	//	fmt.Println("Error getting store business hours:", err)
	//	return
	//}
	//
	//// print the store business hours
	//fmt.Println(storeBusinessHour)

	endTime := "2023-01-24 09:07:26.441407 UTC"
	startTime, err := utils.GetTimeOfXDaysBefore(endTime, 7)
	if err != nil {
		return
	}

	localStartTime, _, _ := utils.ConvertUTCStrToLocal(startTime, "America/Denver")
	localEndTime, _, _ := utils.ConvertUTCStrToLocal(endTime, "America/Denver")

	businessHours, err := storeBusinessHourRepo.GetBusinessHoursInTimeRange(context.Background(), int64(storeId), localStartTime, localEndTime)
	if err != nil {
		fmt.Println("Error getting store business hours:", err)
		return
	}
	storeStatus, err := storeStatusRepo.GetStoreStatusInTimeRange(context.Background(), int64(storeId), startTime, endTime)
	if err != nil {
		return
	}
	fmt.Println(storeStatus, businessHours)*/
	app.StartApplication(context.Background())

}
