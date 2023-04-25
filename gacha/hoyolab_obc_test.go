package gacha_test

import (
	"fmt"
	"give-me-genshin-gacha/gacha"
	"testing"
)

func TestGEt(t *testing.T) {
	a, w, err := gacha.GetGameItem("zh-cn")
	if err != nil {
		panic(err)
	}
	for _, v := range a {
		fmt.Println(v.Name)
	}
	for _, v := range w {
		fmt.Println(v.Name)
	}
}
