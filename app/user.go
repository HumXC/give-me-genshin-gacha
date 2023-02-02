package app

import (
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type UserMan struct {
	db      *models.UserDB
	webview *webview.WebView
}

func (u *UserMan) Get() []models.User {
	result, err := u.db.Get()
	if err != nil {
		u.webview.Alert.Error("请求失败: " + err.Error())
		return result
	}
	return result
}

func (u *UserMan) Sync(id uint64, rawUrl string) bool {
	err := u.db.Sync(id, rawUrl)
	if err != nil {
		u.webview.Alert.Error("更新用户信息失败: " + err.Error())
		return false
	}
	return true
}
