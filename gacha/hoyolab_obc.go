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
	ItemType models.ItemType
}

const ItemInfoAPI = "https://sg-public-api-static.hoyolab.com/common/map_user/ys_obc/v1/map/game_item?app_sn=ys_obc"

var resps map[string][]ObcListItem = make(map[string][]ObcListItem)

type AssetsStore struct{}

// 此处的 LoadedLang 与 model.ItemDB 中的 LangedLang 是不一样的
// 这里指的是已经加载的从网络上请求来的热数据
func (a *AssetsStore) LoadedLang() []string {
	keys := make([]string, 0, len(resps))
	for k := range resps {
		keys = append(keys, k)
	}
	return keys
}
func (a *AssetsStore) Get(lang string) ([]ObcListItem, error) {
	_, ok := resps[lang]
	if !ok {
		result, err := a.fetch(lang)
		if err != nil {
			return nil, err
		}
		resps[lang] = result
	}
	return resps[lang], nil
}

func (a *AssetsStore) fetch(lang string) ([]ObcListItem, error) {
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
		return nil, err
	}
	query := _url.Query()
	query.Add("lang", lang)
	_url.RawQuery = query.Encode()
	r, err := http.DefaultClient.Get(_url.String())
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, errors.New(strconv.Itoa(r.StatusCode) + ": " + r.Status)
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	resp := ObcRespBody{}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Message != "OK" {
		return nil, errors.New(resp.Message)
	}
	aList := resp.Data.Avatar.List
	wList := resp.Data.Weapon.List
	for i := 0; i < len(aList); i++ {
		aList[i].ItemType = models.ItemAvatar
	}
	for i := 0; i < len(wList); i++ {
		wList[i].ItemType = models.ItemWeapon
	}
	return append(aList, wList...), nil
}
