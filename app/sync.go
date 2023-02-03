package app

import (
	"context"
	"errors"
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
	"strconv"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SyncMsg struct {
	Uid    string `json:"uid"`
	RawURL string `json:"raw_url"`
	Error  string `json:"error"`
}

// 提供同步祈愿数据的功能
type SyncMan struct {
	ctx     context.Context
	proxy   *gacha.ProxyServer
	config  *config.Config
	webview *webview.WebView
	logDB   *models.LogDB
	userDB  *models.UserDB
}

func (s *SyncMan) Sync(rawUrl string) uint64 {
	f, err := gacha.NewFetcher(rawUrl)
	if err != nil {
		// 测试失败返回空，但是此错误不需要在前端报告
		if errors.Is(err, gacha.ErrUrlTestFailed) {
			s.webview.Alert.Warning("爬虫创建失败: " + err.Error())
			return 0
		}
		if err != nil {
			s.webview.Alert.Error("请求失败: " + err.Error())
			return 0
		}
	}
	gachaTypes := []string{
		"301", "302", "200", "100",
	}
	endIDs, err := s.logDB.EndLogIDs()
	if err != nil {
		s.webview.Alert.Error("获取 EndLogID 失败: " + err.Error())
		return 0
	}
	result := make([]models.GachaLog, 0)

	for _, gachaType := range gachaTypes {
		g := f.Get(gachaType, endIDs[gachaType])
		page := 0
		var err error
		for {
			page++
			runtime.LogInfo(s.ctx, "正在获取["+gachaType+"]第 "+strconv.Itoa(page)+" 页")
			err = nil
			resp, e := g()
			logs := gacha.ConverToDBLog(resp)
			// 保证较新的记录在数组结尾
			for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
				logs[i], logs[j] = logs[j], logs[i]
			}
			result = append(logs, result...)
			if e != nil {
				err = e
				break
			}

		}
		if !errors.Is(err, gacha.ErrPageEnd) {
			s.webview.Alert.Error("获取祈愿记录失败: " + err.Error())
		} else {
			runtime.LogInfo(s.ctx, "已经是最后一页了")
		}
	}
	runtime.LogInfo(s.ctx, "将祈愿记录写入数据库")
	err = s.logDB.Add(result)
	if err != nil {
		s.webview.Alert.Error("祈愿记录写入数据库失败: " + err.Error())
		return 0
	}

	return f.Uid()
}
func (s *SyncMan) StopProxyServer() bool {
	if s.proxy == nil {
		return true
	}
	err := s.proxy.Stop()
	if err != nil {
		s.webview.Alert.Error("无法关闭代理服务器: " + err.Error())
	}
	runtime.LogInfo(s.ctx, "手动动关闭代理服务器")
	return err == nil
}
func (s *SyncMan) GetRawURL(isUseProxy bool) string {
	if isUseProxy {
		if s.proxy == nil {
			runtime.LogInfo(s.ctx, "创建代理服务器")
			p, err := gacha.NewProxyServer()
			if err != nil {
				s.webview.Alert.Error("无法创建代理服务器: " + err.Error())
				return ""
			}
			s.proxy = p
		}
		runtime.LogInfo(s.ctx, "开始从代理服务器获取链接")
		url, err := s.proxy.Start(gacha.Api)
		if err != nil {
			s.webview.Alert.Error("无法启动理服务器: " + err.Error())
			return ""
		}
		return url
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
		s.webview.Alert.Error("无法获取祈愿链接: " + err.Error())
		return ""
	}
	return rawURL
}
