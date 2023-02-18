package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/store_monitoring/database"
)

var (
	Router = gin.Default()
)

func StartApplication(ctx context.Context) {
	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	SetupRoutes(db)
	Router.Run(":8080")
}
