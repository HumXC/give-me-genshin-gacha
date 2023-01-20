package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uid      string
	RawURL   string
	SyncTime time.Time
}

type UserDB struct {
	uids []string
}

func (u *UserDB) hasUid(uid string) (bool, error) {
	_uid := ""
	tx := db.Table("users").Select("uid").Where("uid=?", uid).First(&_uid)
	return _uid == uid, tx.Error
}

// 此函数表示某个用户已经进行了一次同步操作
func (u *UserDB) Sync(uid, rawURL string) error {
	user := User{
		Uid:      uid,
		RawURL:   rawURL,
		SyncTime: time.Now(),
	}
	hasUid, err := u.hasUid(uid)
	if err != nil {
		return err
	}
	if !hasUid {
		tx := db.Table("users").Create(&user)
		if tx.Error == nil {
			u.uids = append(u.uids, uid)
		}
		return tx.Error
	}
	tx := db.Table("users").Omit("uid").Where("uid=?", uid).Updates(&user)
	return tx.Error
}

func (u *UserDB) Get() ([]User, error) {
	users := make([]User, 4)
	tx := db.Table("users").Find(&users)
	return users, tx.Error
}
