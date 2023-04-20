package gacha

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// 此文件用于请求外部的 hoyolab API 来请求获取物品头像等资源
type GameItem struct {
	Icon     string `json:"icon"`
	ItemID   uint   `json:"item_id"`
	RankType int    `json:"level"`
	Name     string `json:"name"`
}

const APIGameItem = "https://sg-public-api-static.hoyolab.com/common/map_user/ys_obc/v1/map/game_item?app_sn=ys_obc"

func GetGameItem(lang string) ([]GameItem, []GameItem, error) {
	type ObcRespBody struct {
		Data struct {
			Avatar struct {
				List []GameItem `json:"list"`
			} `json:"avatar"`
			Weapon struct {
				List []GameItem `json:"list"`
			} `json:"weapon"`
		} `json:"data"`
		Message string `json:"message"`
	}
	_url, err := url.Parse(APIGameItem)
	if err != nil {
		return nil, nil, err
	}
	query := _url.Query()
	query.Add("lang", lang)
	_url.RawQuery = query.Encode()
	r, err := http.DefaultClient.Get(_url.String())
	if err != nil {
		return nil, nil, err
	}
	if r.StatusCode != 200 {
		return nil, nil, errors.New(strconv.Itoa(r.StatusCode) + ": " + r.Status)
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	resp := ObcRespBody{}
	err = json.Unmarshal(b, &resp)
	if err != nil {
		return nil, nil, err
	}
	if resp.Message != "OK" {
		return nil, nil, errors.New(resp.Message)
	}
	return resp.Data.Avatar.List, resp.Data.Weapon.List, nil
}
