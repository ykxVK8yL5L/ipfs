package inc

import (
	"ipfs/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var DBErr error
var err error

// InitDB opens a database and saves the reference to `Database` struct.
func InitDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("config/upload.db"), &gorm.Config{})
	db.AutoMigrate(&model.Upload{})
	DB = db
	return DB
}

// GetDB helps you to get a connection
func GetDB() *gorm.DB {
	return DB
}

// GetDBErr helps you to get the db error
func GetDBErr() error {
	return DBErr
}
