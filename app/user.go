package app

import (
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type UserMan struct {
	db      *models.UserDB
	webview *webview.WebView
}

func NewUserMan(webview webview.WebView) *UserMan {
	return &UserMan{
		webview: &webview,
		db:      &models.UserDB{},
	}
}
