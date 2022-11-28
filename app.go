package main

import (
	"context"
	"encoding/json"
	"fmt"
	"give-me-genshin-gacha/database"
	"io"
	"os"
	"path"

	"github.com/wailsapp/wails/v2/pkg/logger"
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
} // end

type Bridge struct {
	c chan (Message)
}

// 发送一条消息，一般由后端调用
func (b *Bridge) PutMsg(m Message) {
	b.c <- m
}

// 阻塞接受消息，一般由前端调用
func (b *Bridge) GetMsg() Message {
	m := <-b.c
	return m
}

type App struct {
	ctx     context.Context
	l       logger.Logger
	DB      database.GachaDB
	Bridge  *Bridge
	DataDir string
	GameDir string
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
		m := fmt.Sprintf("加载配置文件时出现错误: %s", err.Error())
		a.Bridge.PutMsg(Message{
			Type: "error",
			Msg:  m,
		})
		a.l.Error(m)
		return opt
	}
	err = json.Unmarshal(d, &opt)
	if err != nil {
		m := fmt.Sprintf("加载配置文件时出现错误: %s", err.Error())
		a.Bridge.PutMsg(Message{
			Type: "error",
			Msg:  m,
		})
		a.l.Error(m)
	}
	return opt
}
func (a *App) SaveOption(opt Option) {
	name := path.Join(a.DataDir, "option.json")
	b, err := json.Marshal(opt)
	if err != nil {
		m := fmt.Sprintf("保存配置文件时出现错误: %s", err.Error())
		a.Bridge.PutMsg(Message{
			Type: "error",
			Msg:  m,
		})
		a.l.Error(m)
		return
	}
	err = os.Remove(name)
	if err != nil {
		if !os.IsNotExist(err) {
			m := fmt.Sprintf("保存配置文件时出现错误: %s", err.Error())
			a.Bridge.PutMsg(Message{
				Type: "error",
				Msg:  m,
			})
			a.l.Error(m)
			return
		}
	}
	f, err := os.Create(name)
	if err != nil {
		m := fmt.Sprintf("保存配置文件时出现错误: %s", err.Error())
		a.Bridge.PutMsg(Message{
			Type: "error",
			Msg:  m,
		})
		a.l.Error(m)
		return
	}
	defer f.Close()
	f.Write(b)
}
func (a *App) GetPieDatas() {

}
func (a *App) GetUids() []string {
	r, err := a.DB.GetUids()
	if err != nil {
		a.Bridge.PutMsg(Message{
			Type: "error",
			Msg:  "无法获取 uid: " + err.Error(),
		})
		return nil
	}
	a.l.Info(fmt.Sprintf("获取到数据库中存在的 uid: %v", r))
	return r
}

// NewApp creates a new App application struct
func NewApp() *App {
	execP, _ := os.Executable()
	l := logger.NewDefaultLogger()
	dataDir := path.Dir(execP)
	dbName := path.Join(dataDir, "data.db")
	db, err := database.NewDB(dbName)
	if err != nil {
		l.Error(err.Error())
	}
	l.Info("使用了数据库: " + dbName)
	return &App{
		DB:      db,
		l:       l,
		DataDir: dataDir,
		Bridge: &Bridge{
			c: make(chan Message),
		},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
