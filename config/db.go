package config

import (
	"ExchangeApp/global"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDB() {
	//data source name
	dsn := AppConfig.DataBase.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize DataBase,got error: %v", err)
	}

	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(AppConfig.DataBase.MaxIdleConns)

	sqlDB.SetMaxOpenConns(AppConfig.DataBase.MaxOpenConns)

	sqlDB.SetConnMaxIdleTime(time.Hour)

	if err != nil {
		log.Fatalf("Failed to configue DataBase,got error: %v", err)
	}
	global.Db = db
}
