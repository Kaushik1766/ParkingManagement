package main

import (
	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/Kaushik1766/ParkingManagement/internal/models"
)

func main() {
	gormDB, _ := db.InitDB()

	err := db.MigrateModels(
		gormDB,
		models.Building{},
		models.Floor{},
		models.Slot{},
		models.Office{},
		models.User{},
		models.Vehicle{},
		models.ParkingHistory{},
	)
	if err != nil {
		panic("error migrating models" + err.Error())
	}
}
