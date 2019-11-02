package models

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB //database

func init() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password) //Build connection string
	log.Trace(dbUri)

	conn, err := gorm.Open("postgres", dbUri)

	for err != nil {
		conn, err = gorm.Open("postgres", dbUri)

		if err != nil {
			log.WithFields(log.Fields{
				"host": dbHost,
				"user": username,
				"name": dbName,
			}).Warn("Waiting for database to become available")
			log.Debug(err)
		}

		time.Sleep(time.Duration(3) * time.Second)
	}

	log.Info("Connected to database")

	db = conn
	db.Debug().AutoMigrate(&Update{}, &Host{}, &Account{}) //Database migration
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return db
}
