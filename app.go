package main

import (
	"context"
	"encoding/json"
	"fmt"
	"give-me-genshin-gacha/database"
	"give-me-genshin-gacha/gacha"
	"io"
	"os"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// 前端定义和使用的数据, 应该与下面文件里的定义相同
// frontend/src/type.ts
// start
type PieData struct {
	UsedCost   int    `json:"usedCost"`
	Arms3Total int    `json:"arms3Total"`
	Arms4Total int    `json:"arms4Total"`
	Arms5Total int    `json:"arms5Total"`
	Role4Total int    `json:"role4Total"`
	Role5Total int    `json:"role5Total"`
	GachaType  string `json:"gachaType"`
}
type OtherOption struct {
	AutoSync  bool `json:"autoSync"`
	UseProxy  bool `json:"useProxy"`
	DarkTheme bool `json:"darkTheme"`
}
type ShowGacha struct {
	RoleUp    bool `json:"roleUp"`
	ArmsUp    bool `json:"armsUp"`
	Permanent bool `json:"permanent"`
	Start     bool `json:"start"`
}
type ControlBar struct {
	SelectedUid string `json:"selectedUid"`
}
type Option struct {
	ShowGacha   ShowGacha   `json:"showGacha"`
	OtherOption OtherOption `json:"otherOption"`
	ControlBar  ControlBar  `json:"controlBar"`
}
type Message struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

type GachaPieTotals struct {
	T301 []database.GachaTotal `json:"t301"`
	T302 []database.GachaTotal `json:"t302"`
	T200 []database.GachaTotal `json:"t200"`
	T100 []database.GachaTotal `json:"t100"`
}
type GachaPieDate struct {
	UsedCosts []database.GachaUsedCost `json:"usedCosts"`
	Totals    GachaPieTotals           `json:"totals"`
}

type App struct {
	ctx     context.Context
	proxy   *gacha.ProxyServer
	DB      *database.GachaDB
	DataDir string
	GameDir string
}

func (a *App) putErr(info string, err error) {
	m := fmt.Sprint(info, "-", err)
	runtime.EventsEmit(a.ctx, "alert", Message{
		Type: "error",
		Msg:  m,
	})
	runtime.LogError(a.ctx, m)
}
func (a *App) GetOption() Option {
	name := path.Join(a.DataDir, "option.json")
	opt := Option{
		ShowGacha: ShowGacha{
			RoleUp: true,
		},
	}
	f, err := os.Open(name)
	if os.IsNotExist(err) {
		a.SaveOption(opt)
		return opt
	}
	defer f.Close()
	d, err := io.ReadAll(f)
	if err != nil {
		a.putErr("加载配置文件时出现错误", err)
		return opt
	}
	err = json.Unmarshal(d, &opt)
	if err != nil {
		a.putErr("加载配置文件时出现错误", err)
	}
	return opt
}
func (a *App) SaveOption(opt Option) {
	name := path.Join(a.DataDir, "option.json")
	b, err := json.Marshal(opt)
	if err != nil {
		a.putErr("保存配置文件时出现错误", err)
		return
	}
	err = os.Remove(name)
	if err != nil {
		if !os.IsNotExist(err) {
			a.putErr("保存配置文件时出现错误", err)
			return
		}
	}
	f, err := os.Create(name)
	if err != nil {
		a.putErr("保存配置文件时出现错误", err)
		return
	}
	defer f.Close()
	f.Write(b)
}

// 饼图数据
func (a *App) GetPieDatas(uid string) GachaPieDate {
	result := GachaPieDate{}
	r, err := a.DB.GetTotals(uid)
	if err != nil {
		return result
	}
	c, err := a.DB.GetUsedCost(uid)
	if err != nil {
		return result
	}
	result.UsedCosts = c
	result.Totals.T301 = r["301"]
	result.Totals.T302 = r["302"]
	result.Totals.T200 = r["200"]
	result.Totals.T100 = r["100"]

	return result
}
func (a *App) GetUids() []string {
	return a.DB.Uids
}

// 从服务器同步祈愿数据到本地数据库, 如果成功返回 true
func (a *App) Sync(useProxy bool) string {
	rawUrl := ""
	if useProxy {
		proxy, err := gacha.NewProxyServer()
		if err != nil {
			a.putErr("代理服务器创建失败", err)
			return "fail"
		}
		a.proxy = proxy
		err = proxy.Start()
		runtime.EventsEmit(a.ctx, "proxy-started")
		if err != nil {
			a.putErr("代理服务器启动失败", err)
			return "fail"
		}
		rawUrl = <-proxy.Url
		a.proxy = nil
		if rawUrl == "" {
			return "cancel"
		}
		err = proxy.Stop()
		runtime.EventsEmit(a.ctx, "proxy-stoped")
		if err != nil {
			a.putErr("代理服务器关闭失败", err)
			return "fail"
		}
	} else {
		if a.GameDir == "" {
			dir, err := gacha.GetGameDir()
			if err != nil {
				a.putErr("获取游戏目录失败", err)
				return "fail"
			}
			a.GameDir = dir
		}
		url, err := gacha.GetRawURL(a.GameDir)
		if err != nil {
			a.putErr("无法获取祈愿链接", err)
			return "fail"
		}
		rawUrl = url
	}

	fetcher, err := gacha.NewFetcher(a.ctx, rawUrl)
	if err != nil {
		a.putErr("无法创建爬虫", err)
		return "fail"
	}
	lastIds, err := a.DB.GetLastIDs()
	if err != nil {
		a.putErr("无法从数据库获取最新的物品", err)
		return "fail"
	}
	items, err := fetcher.Get(lastIds)
	if err != nil {
		if err.Error() == "authkey timeout" {
			return "authkey timeout"
		}
		runtime.EventsEmit(a.ctx, "alert", Message{
			Type: "warning",
			Msg:  "从服务器获取数据时出现错误, 可能无法同步所有数据 - " + err.Error(),
		})
	}
	err = a.DB.Add(fetcher.Uid, items)
	if err != nil {
		a.putErr("写入数据库失败", err)
		return "fail"
	}
	return fetcher.Uid
}

// NewApp creates a new App application struct
func NewApp() *App {
	execP, _ := os.Executable()
	dataDir := path.Dir(execP)
	dbName := path.Join(dataDir, "data.db")
	db, err := database.NewDB(dbName)
	if err != nil {
		panic(err)
	}
	return &App{
		DB:      db,
		DataDir: dataDir,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 前端发来的关闭代理请求
	runtime.EventsOn(a.ctx, "stop-proxy", func(optionalData ...interface{}) {
		if a.proxy == nil {
			return
		}
		err := a.proxy.Stop()
		if err != nil {
			a.putErr("关闭代理时出现错误", err)
		}
	})
}

func (a *App) shutdown(ctx context.Context) {
	if a.proxy != nil {
		_ = a.proxy.Stop()
	}
}
