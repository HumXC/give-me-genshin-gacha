package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// 初始化数据库
func InitDB(name string) error {
	if db != nil {
		return nil
	}
	d, err := gorm.Open(sqlite.Open(name), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}
	d.AutoMigrate(&GachaLog{}, &User{})
	db = d
	return nil
}
