package app

import (
	"context"
	"errors"
	"give-me-genshin-gacha/assets"
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var ErrLangSrcOutdate = errors.New("lang src outdate")

type SyncMsg struct {
	Uid    string `json:"uid"`
	RawURL string `json:"raw_url"`
	Error  string `json:"error"`
}

// 提供同步祈愿数据的功能
type SyncMan struct {
	ctx       context.Context
	proxy     *gacha.ProxyServer
	config    *config.Config
	webview   *webview.WebView
	logDB     *models.LogDB
	userDB    *models.UserDB
	itemStore *assets.ItemStore
	itemDB    *models.ItemDB
}

func (s *SyncMan) Sync(rawURL string) uint64 {
	f, err := gacha.NewFetcher(rawURL)
	if err != nil {
		if errors.Is(err, gacha.ErrURLAuthTimeout) {
			s.webview.Alert.Warning("URL 已过期，请尝试")
			return 0
		}
		if errors.Is(err, gacha.ErrURLTestFailed) {
			s.webview.Alert.Error("爬虫创建失败: " + err.Error())
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
		get := f.Get(gachaType, endIDs[gachaType])
		page := 0
		isPageEnd := false
		for {
			page++
			runtime.LogInfo(s.ctx, "正在获取["+gachaType+"]第 "+strconv.Itoa(page)+" 页")
			resp, err := get()
			if err != nil {
				if !errors.Is(err, gacha.ErrPageEnd) {
					s.webview.Alert.Error("获取祈愿记录失败: " + err.Error())
					return 0
				} else {
					isPageEnd = true
					runtime.LogInfo(s.ctx, "已经是最后一页了")
				}
			}
			// 加载对应的语言资源，否则无法找到对应的 ItemID
			_err := s.itemStore.Load(f.Lang())
			if _err != nil {
				s.webview.Alert.Error("加载物品信息失败: " + err.Error())
				return 0
			}
			logs, _err := s.converToDBLog(resp)
			if _err != nil {
				if errors.Is(_err, ErrLangSrcOutdate) {
					s.webview.Alert.Warning("物品信息已过期，请重试以更新: " + _err.Error())
					return 0
				}
				s.webview.Alert.Error("数据转换失败: " + _err.Error())
				return 0
			}
			// 保证较新的记录在数组结尾
			for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
				logs[i], logs[j] = logs[j], logs[i]
			}
			result = append(logs, result...)
			if isPageEnd {
				break
			}
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
	runtime.LogInfo(s.ctx, "手动关闭代理服务器")
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
			s.webview.Alert.Error("无法获取游戏目录" + err.Error())
			return ""
		}
		s.config.GameDir = dir
	}
	rawURL, err := gacha.GetRawURL(s.config.GameDir)
	if err != nil {
		s.webview.Alert.Error("无法获取祈愿链接: " + err.Error())
		// 获取失败也可能是游戏目录的问题，所以清除游戏目录
		s.config.GameDir = ""
		return ""
	}
	return rawURL
}

func (s *SyncMan) converToDBLog(src []gacha.RespDataListItem) ([]models.GachaLog, error) {
	layout := "2006-01-02 15:04:05"
	result := make([]models.GachaLog, 0)
	for i := 0; i < len(src); i++ {
		log := models.GachaLog{}
		log.OriginGachaType = src[i].GachaType
		if log.OriginGachaType == "400" {
			log.GachaType = "301"
		} else {
			log.GachaType = log.OriginGachaType
		}
		// 如果没有找到对应的物品，那么查询出来的结果 id 等于 0，这就说明语言资源过期了
		item, err := s.itemDB.GetWithName(src[i].Lang, src[i].Name)
		if err != nil {
			return nil, err
		}
		if item.ID == 0 {
			s.itemStore.UnLoad(src[i].Lang)
			return nil, ErrLangSrcOutdate
		}
		logID, err := strconv.ParseUint(src[i].ID, 10, 64)
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(src[i].Count)
		if err != nil {
			return nil, err
		}
		uid, err := strconv.ParseUint(src[i].Uid, 10, 64)
		if err != nil {
			return nil, err
		}
		t, err := time.ParseInLocation(layout, src[i].Time, time.Local)
		if err != nil {
			return nil, err
		}
		log.ItemID = item.ID
		log.LogID = logID
		log.Count = count
		log.Uid = uid

		log.Time = t
		result = append(result, log)
	}
	return result, nil
}
