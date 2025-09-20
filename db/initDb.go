package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// err := godotenv.Load(config.EnvPath)
	// if err != nil {
	// 	panic("Error loading .env file")
	// }

	dbURL := os.Getenv("DATABASE_URL")
	return gorm.Open(postgres.Open(dbURL), &gorm.Config{})
}

func MigrateModels(db *gorm.DB, models ...any) error {
	if db == nil {
		return nil
	}
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		log.Fatal("failed to create uuid-ossp extension: ", err)
	}
	//db.Migrator().DropTable(models...)
	// for _, model := range models {
	// 	if err := db.AutoMigrate(model); err != nil {
	// 		return err
	// 	}
	// }
	db.AutoMigrate(models...)
	fmt.Println("models migrated successfully")
	return nil
}
