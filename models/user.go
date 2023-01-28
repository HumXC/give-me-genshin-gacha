package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	RawURL   string
	SyncTime time.Time
}

type UserDB struct{}

// 此函数表示某个用户已经进行了一次同步操作
func (u *UserDB) Sync(id uint, rawURL string) error {
	user := User{
		Model: gorm.Model{
			ID: id,
		},
		RawURL:   rawURL,
		SyncTime: time.Now(),
	}
	return db.Model(&User{}).Save(&user).Error
}

func (u *UserDB) Get() ([]User, error) {
	users := make([]User, 4)
	tx := db.Model(&User{}).Find(&users)
	return users, tx.Error
}
