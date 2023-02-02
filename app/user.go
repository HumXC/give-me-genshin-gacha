package app

import (
	"give-me-genshin-gacha/models"
)

type UserMan struct {
	db *models.UserDB
}

func (u *UserMan) Get() []models.User {
	result, err := u.db.Get()
	if err != nil {
		// TODO: 处理错误
		return result
	}
	return result
}

func (u *UserMan) Sync(id uint, rawUrl string) {
	err := u.db.Sync(id, rawUrl)
	if err != nil {
		// TODO: 处理错误
		return
	}
}
