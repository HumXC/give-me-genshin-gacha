package app

import (
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type SyncMsg struct {
	Uid    string `json:"uid"`
	RawURL string `json:"raw_url"`
	Error  string `json:"error"`
}

// 提供同步祈愿数据的功能
type SyncMan struct {
	Proxy   *gacha.ProxyServer
	config  *config.Config
	userDB  models.UserDB
	webview *webview.WebView
}

// 同步有两种情况
// 1、从数据库里的 users 表获取已知的 rawURL 进行同步，同步失败就清除对应的 rawURL
// 2、从游戏目录或者代理服务器中获取 rawURL, 然后获取到 UID, 存入 users 表
// 分成两个动作, 获取rawURL 和 同步数据
