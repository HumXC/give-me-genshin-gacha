package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RawURL    string         `json:"raw_url"`
	SyncTime  time.Time      `json:"sync_time"`
}

type UserDB struct {
	db *gorm.DB
}

// 此函数表示某个用户已经进行了一次同步操作
func (u *UserDB) Sync(id int, rawURL string) error {
	user := User{}
	err := u.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return err
	}
	user.ID = id
	user.SyncTime = time.Now()
	user.RawURL = rawURL
	// 如果 rawURL=="" 说明原来存在数据库里的链接已经失效，则不更新 SyncTime
	if rawURL == "" {
		return u.db.Model(&User{}).Where("id = ?", id).Update("raw_url", "").Error
	}
	return u.db.Save(&user).Error
}

func (u *UserDB) Get() ([]User, error) {
	users := make([]User, 0)
	err := u.db.Model(&User{}).Find(&users).Error
	return users, err
}
