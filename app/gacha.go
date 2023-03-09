package app

import (
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type GachaMan struct {
	logDB   *models.LogDB
	webview *webview.WebView
}

func (g *GachaMan) GetGachaInfo(uid int) []models.GachaInfo {
	r, err := g.logDB.GetInfo(uid)
	if err != nil {
		g.webview.Alert.Error("从数据库获取祈愿信息失败: " + err.Error())
		return nil
	}
	return r
}

// TODO: 优化此处 where
func (g *GachaMan) GetGachaLogs(page int, uid uint, lang, gachaType string, filter config.FilterOption, desc bool) []models.FullGachaLog {
	if lang == "" {
		lang = "zh-cn"
	}
	result, err := g.logDB.GetFullGacha(page, uid, lang, gachaType, filter.Avatar4, filter.Avatar5, filter.Weapon3, filter.Weapon4, filter.Weapon5, desc)
	if err != nil {
		g.webview.Alert.Error("从数据库获取祈愿记录失败: " + err.Error())
		return nil
	}
	return result
}
