package webview

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Message struct {
	Type string `json:"type"`
	Msg  string `json:"message"`
}

type Alert struct{ ctx context.Context }

// 发送 alert
func (m *Alert) Send(message Message) {
	runtime.EventsEmit(m.ctx, "alert", message)
}
func (m *Alert) Error(msg string) {
	m.Send(Message{
		Type: "error",
		Msg:  msg,
	})
}
func (m *Alert) Success(msg string) {
	m.Send(Message{
		Type: "success",
		Msg:  msg,
	})
}
func (m *Alert) Info(msg string) {
	m.Send(Message{
		Type: "info",
		Msg:  msg,
	})
}
func (m *Alert) Warning(msg string) {
	m.Send(Message{
		Type: "warning",
		Msg:  msg,
	})
}
