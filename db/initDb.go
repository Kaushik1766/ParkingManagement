package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDb(connectionUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionUrl), &gorm.Config{})
	return db, err
}
