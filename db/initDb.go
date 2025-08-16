package db

import (
	"fmt"

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
	// db.Migrator().DropTable(models...)
	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}
	fmt.Println("models migrated successfully")
	return nil
}
