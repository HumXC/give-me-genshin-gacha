package gacha

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	// 角色活动祈愿
	GachaType301 = "301"
	// 武器活动祈愿
	GachaType302 = "302"
	// 常驻祈愿
	GachaType200 = "200"
	// 新手祈愿
	GachaType100 = "100"
)

var ErrURLTestFailed = errors.New("url test failed")
var ErrURLAuthTimeout = errors.New("url authkey timeout")
var ErrPageEnd = errors.New("page end")

type RespDataListItem struct {
	Uid       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	ID        string `json:"id"`
	Count     string `json:"count"`
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
type Fetcher struct {
	uid    uint64
	lang   string
	lock   sync.Mutex
	rawURL *url.URL
}

// 获取此祈愿链接的 Uid
func (f *Fetcher) Uid() uint64 {
	return f.uid
}

// 获取此祈愿链接的 Uid
func (f *Fetcher) Lang() string {
	return f.lang
}

// 进行一次请求, 测试 URL 是否可用
func Test(rawURL string) error {
	type response struct {
		Message string `json:"message"`
		Retcode int    `json:"retcode"`
	}
	resp := response{}
	url := rawURL
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}
	if resp.Retcode == -101 {
		return ErrURLAuthTimeout
	}
	if resp.Retcode != 0 {
		return fmt.Errorf("%w: "+resp.Message, ErrURLTestFailed)
	}
	return nil
}

// 返回一个用于获取祈愿数据的闭包
// 如果没有下一页了, 则返回 ErrPageEnd 错误
// endID 是上一次同步所获取的最后一个物品的 id
func (f *Fetcher) Get(gachaType string, endIDs map[uint64]uint64) func() ([]RespDataListItem, error) {
	isEnd := false
	// end 用于翻页
	end := "0"
	return func() ([]RespDataListItem, error) {
		if isEnd {
			return []RespDataListItem{}, ErrPageEnd
		}
		result, err := f.fetch(gachaType, end)
		if err != nil {
			return result, err
		}
		length := len(result)
		if length == 0 {
			return result, ErrPageEnd
		}
		end = result[length-1].ID
		// 筛选出更新的物品
		for i, item := range result {
			if f.uid == 0 {
				i, _ := strconv.ParseUint(item.Uid, 10, 64)
				f.uid = i
			}
			if f.lang == "" {
				f.lang = item.Lang
			}
			if item.ID == strconv.FormatUint(endIDs[f.uid], 10) {
				isEnd = true
				return result[:i], ErrPageEnd
			}
		}
		return result, nil
	}
}

// 获取指定祈愿的所有记录，gachaType 是数字代号的字符串
// page 的值从 1 开始
func (f *Fetcher) fetch(gachaType, endID string) ([]RespDataListItem, error) {
	result := make([]RespDataListItem, 0)
	f.lock.Lock()
	url := f.rawURL
	query := url.Query()
	query.Set("gacha_type", gachaType)
	query.Set("end_id", endID)
	url.RawQuery = query.Encode()
	url_ := url.String()
	resp, err := http.Get(url_)
	if err != nil {
		return result, err
	}
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	var r Response
	err = json.Unmarshal(jsonData, &r)
	if err != nil {
		return result, err
	}
	if r.Message != "OK" {
		return result, errors.New(r.Message)
	}
	l := r.Data.List
	result = append(result, l...)
	time.Sleep(300 * time.Millisecond)
	f.lock.Unlock()
	return result, nil
}

// 根据祈愿的链接创建爬虫
func NewFetcher(rawURL string) (*Fetcher, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("size", "20")
	query.Set("page", "1")
	query.Set("begin_id", "0")
	u.RawQuery = query.Encode()
	f := &Fetcher{
		rawURL: u,
	}
	return f, Test(rawURL)
}
