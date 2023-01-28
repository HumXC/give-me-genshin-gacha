package main

import (
	"embed"
	"give-me-genshin-gacha/app"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	a := app.NewApp()
	// Create application with options
	err := wails.Run(&options.App{
		Title:     "give-me-genshin-gacha",
		Width:     388,
		Height:    518,
		MinWidth:  388,
		MinHeight: 518,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.Startup(a),
		OnShutdown: app.Shutdown(a),
		Bind: []interface{}{
			a,
			a.GachaMan,
			a.SyncMan,
			a.UserMan,
		},
	},
	)
	if err != nil {
		println("Error:", err.Error())
	}
}
