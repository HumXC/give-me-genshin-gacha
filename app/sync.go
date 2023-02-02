package app

import (
	"errors"
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/webview"
)

type SyncMsg struct {
	Uid    string `json:"uid"`
	RawURL string `json:"raw_url"`
	Error  string `json:"error"`
}

// 提供同步祈愿数据的功能
type SyncMan struct {
	proxy   *gacha.ProxyServer
	config  *config.Config
	webview *webview.WebView
}

// 同步有两种情况
// 1、从数据库里的 users 表获取已知的 rawURL 进行同步，同步失败就清除对应的 rawURL
// 2、从游戏目录或者代理服务器中获取 rawURL, 然后获取到 UID, 存入 users 表
// 分成两个动作, 获取rawURL 和 同步数据

func (s *SyncMan) Sync(rawUrl string) uint {
	f, err := gacha.NewFetcher(rawUrl)
	if err != nil {
		// 测试失败返回空，但是此错误不需要在前端报告
		if errors.Is(err, gacha.ErrUrlTestFailed) {
			s.webview.Alert.Warning("旧链接已经失效，尝试重新获取")
			return 0
		}
		if err != nil {
			s.webview.Alert.Error("请求失败: " + err.Error())
			return 0
		}
	}
	// TODO
	return f.Uid()
}

func (s *SyncMan) GetRawURL(isUseProxy bool) string {
	if isUseProxy {
		// TODO
		return ""
	}

	if s.config.GameDir == "" {
		dir, err := gacha.GetGameDir()
		if err != nil {
			s.webview.Alert.Error(err.Error())
			return ""
		}
		s.config.GameDir = dir
	}
	rawURL, err := gacha.GetRawURL(s.config.GameDir)
	if err != nil {
		s.webview.Alert.Error(err.Error())
		return ""
	}
	return rawURL
}
