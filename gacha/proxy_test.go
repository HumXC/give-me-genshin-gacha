package gacha_test

import (
	"context"
	"fmt"
	"give-me-genshin-gacha/gacha"
	"testing"
)

func TestProxy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())

	server := gacha.NewProxyServer(ctx,
		"https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog")
	fmt.Println(server.Url())
	cancel()
	for err := range server.Err {
		fmt.Println(err)
	}
}
