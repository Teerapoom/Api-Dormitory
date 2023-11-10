package database

import (
	"os"

	"github.com/teerapoom/Dormitory_Api/database/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

// var err error

func IntnDb() {
	dsn := os.Getenv("MYSQL_DNS")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&model.User_Register{})
	Db = db
}
