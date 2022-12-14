package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

// var DB = initializers.ConnectToDb()
var DSN = "root:root@123@tcp(127.0.0.1:3306)/report?charset=utf8mb4&parseTime=True&loc=Local"

func ConnectToDb() *gorm.DB {

	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		panic("DB NOT CONNECTED")
	}
	return DB
}
