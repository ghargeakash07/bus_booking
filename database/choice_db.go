package database

import (
	// "infrapm-rj-go/config"
	"akash-project-go/config"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var BookingDB *gorm.DB

func InitBookingDB() error {

	//Read database configuration from your config package or .env file
	dbConfig := config.GetChoiceDBConfig()

	// Construct the connection string
	dsn := dbConfig.Username + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"

	// Open a database connection
	var err error
	BookingDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error("Failed to connect to the database:", err)
		return err
	}

	// Set up connection pool settings
	sqlDB, err := BookingDB.DB()
	if err != nil {
		logrus.Error("Failed to set up 'BookingDB' database connection pool: ", err)
	}

	// Maximum number of idle connections
	sqlDB.SetMaxIdleConns(10)

	// Maximum number of open connections
	sqlDB.SetMaxOpenConns(100)

	// Maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
