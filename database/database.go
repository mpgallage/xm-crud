package database

import (
	"github.com/jinzhu/gorm"
	"github.com/mpgallage/xmcrud/models"
	log "github.com/sirupsen/logrus"
	"os"
)

var _DATABASE *gorm.DB

func Init() {
	if _DATABASE == nil {
		db, err := gorm.Open("postgres", os.Getenv("DATABASE_ARGS"))
		if err != nil {
			log.Fatal("Error connecting to database.", err)
			panic("Failed to connect database")
		}

		tx := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
		if tx.Error != nil {
			log.Fatal("Error creating extension.", tx.Error)
			return
		}

		db.AutoMigrate(&models.Company{})
		db.AutoMigrate(&models.User{})
		_DATABASE = db
	}
}

func Close() {
	err := _DATABASE.Close()
	if err != nil {
		log.Error("Error closing database!", err)
	}
}

func Get() *gorm.DB {
	return _DATABASE
}
