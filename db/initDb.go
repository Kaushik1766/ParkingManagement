package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(connectionURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(connectionURL), &gorm.Config{})
}

func MigrateModels(db *gorm.DB, models ...any) error {
	if db == nil {
		return nil
	}
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	return nil
}
