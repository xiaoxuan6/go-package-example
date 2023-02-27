package common

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"package-example/models"
)

var DB *gorm.DB

func Init()  {
	db, err := gorm.Open(sqlite.Open("sql.db"), &gorm.Config{})
	if err != nil {
		panic("数据库初始化错误")
	}

	DB = db

	dbAutoMigrate()
}

func dbAutoMigrate() {
	_ = DB.AutoMigrate(models.Address{})
}