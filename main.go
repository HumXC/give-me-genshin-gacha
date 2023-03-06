package main

import (
	"embed"
	"give-me-genshin-gacha/app"
	as "give-me-genshin-gacha/assets"
	"give-me-genshin-gacha/config"
	"give-me-genshin-gacha/models"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	db, err := models.NewDB("./data.db")
	if err != nil {
		println("Error:", err.Error())
		return
	}
	cfg, err := config.Get("./config.json")
	if err != nil {
		println("Error:", err.Error())
		return
	}
	// as 是 assets 包，由于命名冲突便重命名成 as
	itemStore, err := as.NewItemStore(db.Item)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	assetsServer, err := as.NewServer()
	if err != nil {
		println("Error:", err.Error())
		return
	}
	a := app.NewApp(cfg, db, itemStore)
	// Create application with options
	err = wails.Run(&options.App{
		Title:  "give-me-genshin-gacha",
		Width:  740,
		Height: 530,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: assetsServer,
		},
		OnStartup:  app.Startup(a),
		OnShutdown: app.Shutdown(a),
		Bind: []interface{}{
			a,
			a.UserMan,
			a.SyncMan,
			a.GachaMan,
		},
	},
	)
	if err != nil {
		println("Error:", err.Error())
	}
}
