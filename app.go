package main

import (
	"context"
	"fmt"
	"give-me-genshin-gacha/database"
	"give-me-genshin-gacha/gacha"
	"os"
	"path"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

// App struct
type App struct {
	ctx         context.Context
	db          database.GachaDB
	l           logger.Logger
	proxyServer *gacha.ProxyServer
	Error       string
	GameDir     string
}

type GachaInfo struct {
	Name    string      `json:"name"`
	Count   int         `json:"count"`
	S3      int         `json:"s3"`
	S4      int         `json:"s4"`
	S5      int         `json:"s5"`
	S5Items []ItemsInfo `json:"s5Items"`
}
type ItemsInfo struct {
	Name string `json:"name"`
	UseC int    `json:"usec"`
	Time int64  `json:"time"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	add, _ := gacha.SystemProxyAddr()
	enable, _ := gacha.IsEnableSystemProxy()
	fmt.Println(add, enable)
	execP, _ := os.Executable()
	l := logger.NewDefaultLogger()
	dbName := path.Join(path.Dir(execP), "data.db")
	db, err := database.NewDB(dbName)
	if err != nil {
		l.Error(err.Error())
	}
	l.Info("使用了数据库: " + dbName)
	return &App{
		db: db,
		l:  l,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) outErr(msg string) {
	a.Error = msg
	a.l.Error(msg)
}

func (a *App) Uids() []string {
	u, err := a.db.GetUids()
	if err != nil {
		a.outErr("从数据库获取 uid 时出现错误: " + err.Error())
		return nil
	}
	return u
}

// 从米哈游同步祈愿记录到本地
func (a *App) Sync() (uid string) {
	if a.GameDir == "" {
		gameDir, err := gacha.GetGameDir()
		if err != nil {
			a.outErr(fmt.Sprintf("获取游戏目录失败: %s", err))
			return
		}
		a.GameDir = gameDir
	}
	url, err := gacha.GetRawURL(a.GameDir)
	if err != nil {
		a.outErr(fmt.Sprintf("解析游戏缓存时失败: %s", err))
		return
	}
	f, err := gacha.NewFetcher(a.l, url)
	if err != nil {
		a.outErr(fmt.Sprintf("创建爬虫失败: %s", err))
		return
	}
	l, err := a.db.GetLastIDs()
	if err != nil {
		a.outErr(fmt.Sprintf("从数据库获取最后的物品时出现错误: %s", err))
		return
	}
	data, err := f.Get(l)
	if err != nil {
		a.outErr(fmt.Sprintf("从服务器获取祈愿记录时出现错误: %s", err))
		return
	}
	err = a.db.Add(*data)
	if err != nil {
		a.outErr(fmt.Sprintf("添加数据到本地数据库时出现错误: %s", err))
		return
	}
	uid = f.Uid
	return
}

// 获取池子信息
// gachaType 是数字字符串
func (a *App) GetGachaInfo(uid, gachaType string) GachaInfo {
	g := GachaInfo{
		S5Items: make([]ItemsInfo, 0),
	}
	g.Name = gacha.ParseGachaType(gachaType)
	count, err := a.db.CountGacha(uid, gachaType)
	if err != nil {
		a.outErr("从数据库获取祈愿记录时出现错误: " + err.Error())
		return g
	}
	g.Count = count
	g.S3, err = a.db.CountGachaRank(uid, gachaType, "3")
	if err != nil {
		a.outErr("从数据库获取祈愿记录时出现错误: " + err.Error())
		return g
	}
	g.S4, err = a.db.CountGachaRank(uid, gachaType, "4")
	if err != nil {
		a.outErr("从数据库获取祈愿记录时出现错误: " + err.Error())
		return g
	}
	items, err := a.db.GetGachaRank(uid, gachaType, "5")
	if err != nil {
		a.outErr("从数据库获取祈愿记录时出现错误: " + err.Error())
		return g
	}
	g.S5 = len(items)

	for i, v := range items {
		info := ItemsInfo{}
		if i < len(items)-1 {
			count, err := a.db.CountIn(uid, gachaType, items[i+1].ID, v.ID)
			if err != nil {
				a.outErr("从数据库统计祈愿数量时出现错误: " + err.Error())
				return g
			}
			info.Name = v.Name
			info.UseC = count
			info.Time = v.Time.Unix()
		} else {
			oldest, err := a.db.GetOldest(uid, gachaType)
			if err != nil {
				a.outErr("从数据库统计祈愿数量时出现错误: " + err.Error())
				return g
			}
			count, err := a.db.CountIn(uid, gachaType, oldest.ID, v.ID)
			if err != nil {
				a.outErr("从数据库统计祈愿数量时出现错误: " + err.Error())
				return g
			}
			info.Name = v.Name
			info.UseC = count
			info.Time = v.Time.Unix()
		}
		g.S5Items = append(g.S5Items, info)

	}
	return g
}

// 开启代理服务器
func (a *App) StartProxyServer() {
	if a.proxyServer.IsRunning {
		return
	}
	err := a.proxyServer.Start()
	if err != nil {
		a.l.Error("代理服务器开启失败: " + err.Error())
	}
	a.l.Info("成功开启代理服务器")
}

// 关闭代理服务器
func (a *App) StopProxyServer() {
	if !a.proxyServer.IsRunning {
		return
	}
	err := a.proxyServer.Stop()
	if err != nil {
		a.l.Error("代理服务器关闭失败: " + err.Error())
	}
	a.l.Info("成功关闭代理服务器")
}
