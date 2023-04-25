package gacha_test

import (
	"fmt"
	"give-me-genshin-gacha/gacha"
	"testing"
)

func TestXxx(t *testing.T) {
	dir, err := gacha.GetGameDir()
	if err != nil {
		panic(err)
	}
	url, err := gacha.GetRawURL(dir)
	if err != nil {
		panic(err)
	}
	f, err := gacha.NewGachaLogFetcher(url, "zh-cn")
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Uid())
	for _, gachaType := range gacha.GachaTypes() {
		f.SetGachaType(gachaType, "0")
		fmt.Println("GachaType: " + gachaType)
		for f.Next() {
			logs, err := f.Logs()
			if err != nil {
				panic(err)
			}
			fmt.Println(len(logs))
			for _, v := range logs {
				fmt.Println(v.Name)
				break
			}
		}
	}

}
