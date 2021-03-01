package entity

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB initializes the database
func InitDB() error {
	println("InitDB called")
	if db != nil {
		return nil
	}

	d, err := gorm.Open(sqlite.Open("./dbdata/data.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = d

	if err := db.AutoMigrate(&Rule{}, &User{}); err != nil {
		panic(err)
	}

	return nil
}

// Cleanup cleans up any db resources
func Cleanup() {
	// does nothing right now
}
