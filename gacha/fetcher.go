package gacha

import (
	"context"
	"encoding/json"
	"errors"
	"give-me-genshin-gacha/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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
	ItemType  string `json:"item_type"`
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
	ctx context.Context
}

// 获取指定祈愿的所有记录，gachaType 是数字代号的字符串
// lastIDs map[uid]map[gachaType][lastID]
func (f *Fecher) getGacha(gachaTypeNum string, lastIDs map[string]map[string]string) ([]RespDataListItem, error) {
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
		runtime.LogInfof(f.ctx, "正在获取 [%s] 第 %d 页", ParseGachaType(gachaTypeNum), page)
		query.Set("page", strconv.Itoa(page))
		query.Set("end_id", endID)
		f.url.RawQuery = query.Encode()
		url_ := f.url.String()
		resp, err := http.DefaultClient.Get(url_)
		if err != nil {
			return list, err
		}
		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			return list, err
		}
		var r Response
		err = json.Unmarshal(jsonData, &r)
		if err != nil {
			return list, err
		}
		if r.Message != "OK" {
			return list, errors.New(r.Message)
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
			if id, ok := lastIDs[v.Uid][v.GachaType]; !ok {
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
		time.Sleep(300 * time.Millisecond)
	}
	return list, nil
}

func (f *Fecher) Get(lastIDs map[string]map[string]string) ([]models.GachaLog, error) {
	result := make([]models.GachaLog, 0)
	for _, t := range GachaType {
		r, err := f.getGacha(t, lastIDs)
		if err != nil {
			return result, err
		}
		for _, item := range r {
			result = append(result, ConvertToDBItem(item, f.Uid))
		}
		time.Sleep(500 * time.Millisecond)
	}
	return result, nil
}

func ConvertToDBItem(i RespDataListItem, uid string) models.GachaLog {
	return models.GachaLog{
		GachaType: i.GachaType,
		Time:      i.Time,
		Name:      i.Name,
		Lang:      i.Lang,
		ItemType:  i.ItemType,
		RankType:  i.RankType,
		ID:        i.ID,
		Uid:       uid,
	}
}

func NewFetcher(ctx context.Context, rawURL string) (*Fecher, error) {
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
		ctx: ctx,
	}, nil
}
