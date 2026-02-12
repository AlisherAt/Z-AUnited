package database

import (
	"log"

	"project/internal/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg config.Config) *gorm.DB {
	var (
		db  *gorm.DB
		err error
	)

	switch cfg.DBDriver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	default:
		db, err = gorm.Open(sqlite.Open(cfg.DSN), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	DB = db
	return db
}
