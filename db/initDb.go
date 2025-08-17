package db

import (
	"fmt"
	"log"

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
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatal("failed to create uuid-ossp extension: ", err)
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
