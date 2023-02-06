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

// 调用此函数代表此用户已经同步过一次了，根据同步的结果
// 如果成功就会有新的 rawURL，失败则说明数据库中已有的 rawURL 已经不可用，便从数据库中清空
func (u *UserMan) Sync(id uint64, rawUrl string) bool {
	err := u.db.Sync(id, rawUrl)
	if err != nil {
		u.webview.Alert.Error("更新用户信息失败: " + err.Error())
		return false
	}
	return true
}
