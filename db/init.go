package db

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(path string) {
	// open with busy_timeout=5000ms and shared cache
	dsn := "file:" + path + "?cache=shared&_busy_timeout=5000"
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	// Pull out the lower‐level *sql.DB to configure pooling
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// Only one open connection at a time for WAL writes
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	// Enable WAL mode and reasonable synchronous setting
	DB.Exec("PRAGMA journal_mode = WAL;")
	DB.Exec("PRAGMA synchronous = NORMAL;")

	// Auto‐migrate all models
	if err := DB.AutoMigrate(&User{}, &Expression{}, &Task{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
}
