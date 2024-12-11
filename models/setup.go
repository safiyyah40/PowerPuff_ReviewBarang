package models

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
	"log"
	
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:@tcp(localhost:3306)/powerpuff_reviewbarang?charset=utf8mb4&parseTime=True&loc=Local"
    database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
	DB = database
	err = DB.AutoMigrate(&Review{})
    if err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
}
