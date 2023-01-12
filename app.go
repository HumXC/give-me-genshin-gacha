package main

import (
	"context"
	"encoding/json"
	"fmt"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

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
	T301 []models.GachaTotal `json:"t301"`
	T302 []models.GachaTotal `json:"t302"`
	T200 []models.GachaTotal `json:"t200"`
	T100 []models.GachaTotal `json:"t100"`
}
type GachaPieDate struct {
	UsedCosts []models.GachaUsedCost `json:"usedCosts"`
	Totals    GachaPieTotals         `json:"totals"`
}

type App struct {
	ctx         context.Context
	proxy       *gacha.ProxyServer
	option      Option
	DB          *models.GachaDB
	DataDir     string
	GameDir     string
	GachaLogUrl string
}

func (a *App) putErr(info string, err error) {
	m := fmt.Sprint(info, "-", err)
	runtime.EventsEmit(a.ctx, "alert", Message{
		Type: "error",
		Msg:  m,
	})
	runtime.LogError(a.ctx, m)
}
func (a *App) GetNumWithLast(uid, gachaType, id string) int {
	result, err := a.DB.GetNumWithLast(uid, gachaType, id)
	if err != nil {
		a.putErr("从数据库获取计数时出现错误", err)
		return result
	}
	return result
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
	a.option = opt
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
	a.option = opt
}

func (a *App) GetLogs(uid, gachaType string, num, page int) []models.GachaLog {
	result, err := a.DB.GetLogs(uid, gachaType, num, page)
	if err != nil {
		a.putErr("从数据库获取记录时出现错误", err)
		return make([]models.GachaLog, 0)
	}
	return result
}

// 饼图数据
func (a *App) GetPieDatas(uid string) GachaPieDate {
	result := GachaPieDate{}
	r, err := a.DB.GetTotals(uid)
	if err != nil {
		a.putErr("从数据库获取记录时出现错误", err)
		return result
	}
	c, err := a.DB.GetUsedCost(uid)
	if err != nil {
		a.putErr("从数据库获取记录时出现错误", err)
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
	result, err := a.DB.GetUids()
	if err != nil {
		a.putErr("无法从数据库获取 uid", err)
	}
	return result
}

// 从服务器同步祈愿数据到本地数据库, 如果成功返回 true
func (a *App) Sync(useProxy bool) string {
	if useProxy {
		if a.proxy == nil {
			proxy, err := gacha.NewProxyServer()
			if err != nil {
				a.putErr("代理服务器创建失败", err)
				return "fail"
			}
			a.proxy = proxy
		}
		url, err := a.proxy.Start("")
		a.GachaLogUrl = url
		runtime.EventsEmit(a.ctx, "proxy-started")
		if err != nil {
			a.putErr("代理服务器启动失败", err)
			return "fail"
		}
		a.proxy = nil
		if a.GachaLogUrl == "" {
			return "cancel"
		}
		runtime.EventsEmit(a.ctx, "proxy-stoped")
		if err != nil {
			a.putErr("代理服务器关闭失败", err)
			return "fail"
		}
		runtime.EventsEmit(a.ctx, "alert", Message{
			Type: "info",
			Msg:  "成功通过代理取得了祈愿链接，正在同步哦",
		})
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
		a.GachaLogUrl = url
	}

	fetcher, err := gacha.NewFetcher(a.ctx, a.GachaLogUrl)
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
			a.GachaLogUrl = ""
			return "authkey timeout"
		}
		runtime.EventsEmit(a.ctx, "alert", Message{
			Type: "warning",
			Msg:  "从服务器获取数据时出现错误, 可能无法同步所有数据 - " + err.Error(),
		})
	}
	err = a.DB.Add(items)
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
	db, err := models.NewDB(dbName)
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

// 动态资产，用于加载物品图标
type Icon struct {
	Url  string `json:"icon"`
	Name string `json:"name"`
}
type FileLoader struct {
	http.Handler
	IconAvatar []Icon
	IconWeapon []Icon
	IconDir    string
	Inited     bool
}

func NewFileLoader() *FileLoader {
	return &FileLoader{
		IconDir: "icons",
	}
}
func (h *FileLoader) Init() error {
	h.Inited = true
	if !IsExist(h.IconDir) {
		os.MkdirAll(h.IconDir, 0755)
	}
	getIcon := func(url, name string) error {
		f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			h.Inited = false
			return err
		}
		defer f.Close()
		r, err := http.DefaultClient.Get(url)
		if err != nil {
			h.Inited = false
			return err
		}
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			h.Inited = false
			return err
		}
		f.Write(b)
		return nil
	}
	api := "https://waf-api-takumi.mihoyo.com/common/map_user/ys_obc/v1/map/game_item?map_id=2&app_sn=ys_obc&lang=zh-cn"
	if len(h.IconAvatar) == 0 || len(h.IconWeapon) == 0 {
		resp, err := http.DefaultClient.Get(api)
		if err != nil {
			h.Inited = false
			return err
		}
		defer resp.Body.Close()
		type JsonResp struct {
			Data struct {
				Avatar struct {
					List []Icon `json:"list"`
				} `json:"avatar"`
				Weapon struct {
					List []Icon `json:"list"`
				} `json:"weapon"`
			} `json:"data"`
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			h.Inited = false
			return err
		}
		jr := JsonResp{}
		err = json.Unmarshal(body, &jr)
		if err != nil {
			h.Inited = false
			return err
		}
		h.IconAvatar = jr.Data.Avatar.List
		h.IconWeapon = jr.Data.Weapon.List
	}

	for _, icon := range h.IconAvatar {
		name := path.Join(h.IconDir, icon.Name+".png")
		if IsExist(name) {
			continue
		}
		err := getIcon(icon.Url, name)
		if err != nil {
			return err
		}
	}
	for _, icon := range h.IconWeapon {
		name := path.Join(h.IconDir, icon.Name+".png")
		if IsExist(name) {
			continue
		}
		err := getIcon(icon.Url, name)
		if err != nil {
			return err
		}
	}
	return nil
}
func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	println("资产服务器请求:", req.URL.Path)
	iconName := strings.TrimPrefix(req.URL.Path, "/icon/")
	if iconName != "" {
		h.handleIcon(iconName, res)
		return
	}
	http.Error(res, "不支持的路径", http.StatusNotFound)
}
func (h *FileLoader) handleIcon(name string, res http.ResponseWriter) {
	if !h.Inited {
		go func() {
			err := h.Init()
			if err != nil {
				http.Error(res, "图标初始化失败", http.StatusInternalServerError)
				return
			}
		}()
	}
	count := 0
	for !h.Inited {
		time.Sleep(5000)
		if count > 3 {
			http.Error(res, "超时", http.StatusRequestTimeout)
			return
		}
		count++
	}
	iconName := path.Join(h.IconDir, name)
	if IsExist(iconName) {
		f, err := os.ReadFile(iconName)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Write(f)
	} else {
		http.Error(res, "找不到图标", http.StatusNotFound)

	}

}

// 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
