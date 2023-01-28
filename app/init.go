package app

// 此包提供 app 的功能实现
import (
	"context"
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/webview"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx      context.Context
	config   *config.Config
	webView  *webview.WebView
	GachaMan *GachaMan
	SyncMan  *SyncMan
	UserMan  *UserMan
}

// 获取配置
func (a *App) GetConfig() config.Config {
	return *a.config
}

// 修改配置
func (a *App) PutConfig(cfg config.Config) {
	a.config.Put(cfg)
}

func NewApp() *App {
	cfg, err := config.Get("./config.json")
	if err != nil {
		panic(err)
	}
	return &App{
		config:   cfg,
		GachaMan: &GachaMan{},
		SyncMan:  &SyncMan{},
		UserMan:  &UserMan{},
	}
}

func Startup(a *App) func(ctx context.Context) {
	return func(ctx context.Context) {
		a.webView = webview.NewWebView(ctx)
		a.GachaMan.webView = a.webView
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
