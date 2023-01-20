package webview

import "context"

// 提供与前端交互的一些功能

type WebView struct {
	Alert Alert
}

func NewWebView(ctx context.Context) *WebView {
	return &WebView{
		Alert: Alert{
			ctx: ctx,
		},
	}
}
