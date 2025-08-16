package main

import (
	"os"
	"os/user"

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
		user.User{},
		models.Office{},
		models.Vehicle{},
		models.Building{},
		models.Floor{},
		models.Slot{},
		models.ParkingHistory{},
	)
	if err != nil {
		panic("Error migrating models: " + err.Error())
	}
}
