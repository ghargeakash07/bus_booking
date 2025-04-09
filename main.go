package main

import (
	"akash-project-go/database"
	"akash-project-go/routes"
	"akash-project-go/utils"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	utils.InitializeLogger()
}
func main() {
	fmt.Println("shree swami samartha...")

	// Initialize the database
	database.InitBookingDB()
	// Defer a function to close the "insurance" database connection
	defer func() {
		if db, err := database.BookingDB.DB(); err == nil {
			logrus.Error("db not init")
			db.Close()
		}
	}()

	//cretate gin routes
	router := routes.SetupRouter()

	// Create a new gocron scheduler
	ist, _ := time.LoadLocation("maharashtra/pune")
	scheduler := gocron.NewScheduler(ist)
	scheduler.StartAsync()

	// Schedule the task to run every 1 minute
	//scheduler.Every(15).Minute().Do(func() {
	//	services.BreakinSchedular()
	//})

	// Start the server
	router.Run(":" + os.Getenv("APP_PORT"))
	scheduler.Stop()
}
