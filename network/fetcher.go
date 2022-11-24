package network

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var GachaType map[string]string

func init() {
	GachaType = make(map[string]string)
	GachaType["角色活动祈愿"] = "301"
	GachaType["武器活动祈愿"] = "302"
	GachaType["常驻祈愿"] = "200"
	GachaType["新手祈愿"] = "100"
}

func ParseGachaType(gachaType string) string {
	switch gachaType {
	case "301":
		return "角色活动祈愿"
	case "302":
		return "武器活动祈愿"
	case "200":
		return "常驻祈愿"
	case "100":
		return "新手祈愿"
	default:
		return "未知祈愿类型: " + gachaType
	}
}

type RespDataListItem struct {
	Uid       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	ItemId    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	Itemtype  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	ID        string `json:"id"`
}

type RespData struct {
	Page  string             `json:"page"`
	Size  string             `json:"size"`
	Total string             `josn:"total"`
	List  []RespDataListItem `json:"list"`
}
type Response struct {
	RetCode int      `json:"retcode"`
	Message string   `json:"message"`
	Data    RespData `json:"data"`
	Region  string   `json:"region"`
}
type Fecher struct {
	Uid string
	url *url.URL
}

func genError(err error) error {
	return errors.New("Fecher Error: \n\t" + err.Error())
}

// 获取指定祈愿的所有记录，gachaType 是数字代号的字符串
// lastIDs map[uid]map[gachaType][lastID]
func (f *Fecher) Get(gachaTypeNum string, lastIDs map[string]map[string]string) (*[]RespDataListItem, error) {
	list := make([]RespDataListItem, 0)
	page := 1
	endID := "0"
	if len(list) > 0 {
		endID = list[len(list)-1].ID
	}
	query := f.url.Query()

	query.Set("gacha_type", gachaTypeNum)
	for {
		flag := false
		log.Printf("正在获取 [%s] 第 %d 页\n", ParseGachaType(gachaTypeNum), page)
		query.Set("page", strconv.Itoa(page))
		query.Set("end_id", endID)
		f.url.RawQuery = query.Encode()
		url_ := f.url.String()
		resp, err := http.DefaultClient.Get(url_)
		if err != nil {
			return nil, genError(err)
		}
		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, genError(err)
		}
		var r Response
		err = json.Unmarshal(jsonData, &r)
		if err != nil {
			return nil, genError(err)
		}
		if r.Message != "OK" {
			return nil, genError(errors.New("Api error: " + r.Message))
		}
		if len(r.Data.List) == 0 {
			break
		}
		page++
		l := r.Data.List
		if _, ok := lastIDs[f.Uid]; !ok {
			lastIDs[f.Uid] = make(map[string]string)
		}

		// 提取出"更新"的条目
		for i, v := range l {
			if f.Uid == "" {
				f.Uid = v.Uid
			}
			if id, ok := lastIDs[v.Uid][gachaTypeNum]; !ok {
				break
			} else if id != v.ID {
				continue
			}

			list = append(list, l[:i]...)
			flag = true
			break
		}

		if flag {
			break
		}
		endID = l[len(l)-1].ID
		list = append(list, l...)
		time.Sleep(1 * time.Second)
	}
	return &list, nil
}

func NewFetcher(rawURL string) (*Fecher, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("size", "20")
	query.Set("end_id", "0")
	u.RawQuery = query.Encode()
	return &Fecher{
		Uid: "",
		url: u,
	}, nil
}
