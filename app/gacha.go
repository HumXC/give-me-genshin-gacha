package app

import (
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type GachaMan struct {
	db      models.LogDB
	webView *webview.WebView
}
