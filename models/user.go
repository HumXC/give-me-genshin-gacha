package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint64 `gorm:"primarykey" json:"id"`
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
func (u *UserDB) Sync(id uint64, rawURL string) error {
	user := User{
		ID:       id,
		RawURL:   rawURL,
		SyncTime: time.Now(),
	}
	// 如果 rawURL=="" 说明原来存在数据库里的链接已经失效，则不更新 SyncTime
	if rawURL == "" {
		return u.db.Model(&User{}).Where("id = ?", id).Update("raw_url", "").Error
	}
	return u.db.Model(&User{}).Omit("created_at").Where("id = ?", id).Save(&user).Error
}

func (u *UserDB) Get() ([]User, error) {
	users := make([]User, 0)
	err := u.db.Model(&User{}).Find(&users).Error
	return users, err
}
