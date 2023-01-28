package app

import (
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type GachaMan struct {
	db      models.GachaDB
	webView *webview.WebView
}
