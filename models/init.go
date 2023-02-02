package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Item *ItemDB
	Log  *LogDB
	User *UserDB
}

// 初始化数据库
func NewDB(name string) (*DB, error) {
	d, err := gorm.Open(sqlite.Open(name), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	d.AutoMigrate(&GachaLog{}, &User{})
	return &DB{
		Item: &ItemDB{db: d},
		User: &UserDB{db: d},
		Log:  &LogDB{db: d},
	}, nil
}
