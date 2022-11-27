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

type Message struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}
type Bridge struct {
	c chan (Message)
}

// App struct
type App struct {
	ctx     context.Context
	db      database.GachaDB
	l       logger.Logger
	Bridge  *Bridge
	GameDir string
}

// 发送一条消息，只允许后端调用
func (b *Bridge) PutMsg(m Message) {
	b.c <- m
}

// 阻塞接受消息，只允许前端调用
func (b *Bridge) GetMsg() Message {
	m := <-b.c
	return m
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
