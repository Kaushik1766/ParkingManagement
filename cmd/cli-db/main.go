package main

import (
	"os"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
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
		models.Building{},
		models.Floor{},
		models.Slot{},
		models.Office{},
		models.Vehicle{},
		models.User{},
		models.ParkingHistory{},
	)
	if err != nil {
		panic("Error migrating models: " + err.Error())
	}
}
