package entity

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// InitDB initializes the database
func InitDB() error {
	println("InitDB called")
	if db != nil {
		return nil
	}

	d, err := gorm.Open(sqlite.Open("./dbdata/data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	db = d

	if err := db.AutoMigrate(
		&User{},
		&Rule{},
		&Proposal{},
	); err != nil {
		panic(err)
	}

	return nil
}

// Cleanup cleans up any db resources
func Cleanup() {
	// does nothing right now
}
