package utils

import (
	"os"

	"github.com/Kaushik1766/ParkingManagement/db"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func GetDB() (*gorm.DB, error) {
	err := godotenv.Load("/home/kaushik/ParkingManagement/.env")
	if err != nil {
		return nil, err
	}

	dbURL := os.Getenv("DATABASE_URL")
	gormDB, err := db.InitDB(dbURL)
	if err != nil {
		return nil, err
	}

	// err = db.MigrateModels(
	// 	gormDB,
	// 	models.Building{},
	// 	models.Floor{},
	// 	models.Office{},
	// 	models.Slot{},
	// 	models.Vehicle{},
	// 	models.User{},
	// 	models.ParkingHistory{},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	return gormDB, nil
}
