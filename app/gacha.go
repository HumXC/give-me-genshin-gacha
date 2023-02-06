package app

import (
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type GachaMan struct {
	logDB   *models.LogDB
	webview *webview.WebView
}

func (g *GachaMan) GetGachaInfo() []models.GachaInfo {
	r, err := g.logDB.GetInfo()
	if err != nil {
		g.webview.Alert.Error("从数据库获取祈愿信息失败: " + err.Error())
		return nil
	}
	return r
}
