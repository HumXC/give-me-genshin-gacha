package gacha_test

import (
	"fmt"
	"give-me-genshin-gacha/gacha"
	"testing"
	"time"
)

func TestProxy(t *testing.T) {
	server, err := gacha.NewProxyServer()
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(3 * time.Second)
		server.Close()
	}()
	url, err := server.Start("")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)

	server.Close()
	server.Stop()

}
