package main

import (
	"os"
	"os/user"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/models/building"
	"github.com/Kaushik1766/ParkingManagement/internal/models/floor"
	"github.com/Kaushik1766/ParkingManagement/internal/models/office"
	parkinghistory "github.com/Kaushik1766/ParkingManagement/internal/models/parking_history"
	"github.com/Kaushik1766/ParkingManagement/internal/models/slot"
	"github.com/Kaushik1766/ParkingManagement/internal/models/vehicle"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")

	gormDB, err := db.InitDB(dbURL)
	if err != nil {
		panic("Error connecting to the database: " + err.Error())
	}

	err = db.MigrateModels(
		gormDB,
		user.User{},
		office.Office{},
		vehicle.Vehicle{},
		building.Building{},
		floor.Floor{},
		slot.Slot{},
		parkinghistory.ParkingHistory{},
	)
	if err != nil {
		panic("Error migrating models: " + err.Error())
	}
}
