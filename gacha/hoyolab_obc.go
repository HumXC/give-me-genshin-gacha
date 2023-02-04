package gacha

import (
	"encoding/json"
	"errors"
	"give-me-genshin-gacha/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// 此文件用于请求外部的 hoyolab API 来请求获取物品头像等资源
type ObcListItem struct {
	Icon     string `json:"icon"`
	ItemID   int    `json:"item_id"`
	RankType int    `json:"level"`
	Name     string `json:"name"`
	ItemType int
}

const ItemInfoAPI = "https://sg-public-api-static.hoyolab.com/common/map_user/ys_obc/v1/map/game_item?app_sn=ys_obc"

var resps map[string][]ObcListItem = make(map[string][]ObcListItem)

type AssetsStore struct{}

func (a *AssetsStore) Get(lang string) ([]ObcListItem, error) {
	_, ok := resps[lang]
	if !ok {
		a.fetch(lang)
	}
	return resps[lang], nil
}

func (a *AssetsStore) fetch(lang string) error {
	type ObcRespBody struct {
		Data struct {
			Avatar struct {
				List []ObcListItem `json:"list"`
			} `json:"avatar"`
			Weapon struct {
				List []ObcListItem `json:"list"`
			} `json:"weapon"`
		} `json:"data"`
		Message string `json:"message"`
	}
	_url, err := url.Parse(ItemInfoAPI)
	if err != nil {
		return err
	}
	query := _url.Query()
	query.Add("lang", lang)
	_url.RawQuery = query.Encode()
	r, err := http.DefaultClient.Get(_url.String())
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return errors.New(strconv.Itoa(r.StatusCode) + ": " + r.Status)
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	resp := ObcRespBody{}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return err
	}
	if resp.Message != "OK" {
		return errors.New(resp.Message)
	}
	aList := resp.Data.Avatar.List
	wList := resp.Data.Weapon.List
	for i := 0; i < len(aList); i++ {
		aList[i].ItemType = models.ItemAvatar
	}
	for i := 0; i < len(wList); i++ {
		wList[i].ItemType = models.ItemWeapon
	}
	resps[lang] = append(aList, wList...)
	return nil
}
