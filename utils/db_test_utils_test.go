package utils

import (
	"testing"

	"github.com/Kaushik1766/ParkingManagement/db"
	"gorm.io/gorm"
)

func PutDsnInEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://postgres:123@localhost:5432/kaushik")
}

func GetDB() (*gorm.DB, error) {
	gormDB, err := db.InitDB()
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
