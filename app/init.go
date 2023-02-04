package app

// 此包提供 app 的功能实现
import (
	"context"
	"give-me-genshin-gacha/assets"
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	config  *config.Config
	webView *webview.WebView
	UserMan *UserMan
	SyncMan *SyncMan
}

// 获取配置
func (a *App) GetConfig() config.Config {
	return *a.config
}

// 修改配置
func (a *App) PutConfig(cfg config.Config) {
	a.config.Put(cfg)
}

func NewApp(config *config.Config, db *models.DB, itemStore *assets.ItemStore) *App {
	app := &App{
		config:  config,
		UserMan: &UserMan{db: db.User},
		SyncMan: &SyncMan{
			config:    config,
			logDB:     db.Log,
			userDB:    db.User,
			itemStore: itemStore,
			itemDB:    db.Item,
		},
	}
	return app
}

func Startup(a *App) func(ctx context.Context) {
	return func(ctx context.Context) {
		a.ctx = ctx
		webview := webview.NewWebView(ctx)
		a.webView = webview
		a.UserMan.webview = webview
		a.SyncMan.webview = webview
		a.SyncMan.ctx = ctx
	}
}

func Shutdown(a *App) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := a.config.Save()
		if err != nil {
			runtime.LogError(ctx, "配置文件保存失败: "+err.Error())
		}
	}
}
