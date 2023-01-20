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

var ErrUrlTestFailed = errors.New("url test failed")
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
	uid    string
	lock   sync.Mutex
	rawURL *url.URL
}

// 获取此祈愿链接的 Uid
func (f *Fetcher) Uid() string {
	return f.uid
}

// 进行一次请求, 测试 URL 是否可用
func (f *Fetcher) test() error {
	type response struct {
		Message string `json:"message"`
	}
	resp := response{}
	url := f.rawURL.String()
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
	if resp.Message != "OK" {
		return fmt.Errorf(url+" - "+resp.Message, ErrUrlTestFailed)
	}
	return nil
}

// 返回一个用于获取祈愿数据的闭包
// 如果没有下一页了, 则返回 ErrPageEnd 错误
func (f *Fetcher) Get(gachaType, endID string) func() ([]RespDataListItem, error) {
	page := 1
	isEnd := false
	// end 用于翻页
	end := "0"
	return func() ([]RespDataListItem, error) {
		if isEnd {
			return []RespDataListItem{}, ErrPageEnd
		}
		result, err := f.fetch(gachaType, end, page)
		if err != nil {
			return result, err
		}
		page++
		length := len(result)
		end = result[length-1].ID
		// 筛选出更新的物品
		for i, item := range result {
			if f.uid == "" {
				f.uid = item.Uid
			}
			if item.ID == endID {
				isEnd = true
				return result[:i], ErrPageEnd
			}
		}
		return result, nil
	}
}

// 获取指定祈愿的所有记录，gachaType 是数字代号的字符串
// page 的值从 1 开始
func (f *Fetcher) fetch(gachaType, endID string, page int) ([]RespDataListItem, error) {
	result := make([]RespDataListItem, 20)
	f.lock.Lock()
	url := f.rawURL
	query := url.Query()
	query.Set("gacha_type", gachaType)
	query.Set("end_id", endID)
	query.Set("page", strconv.Itoa(page))
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
	u.RawQuery = query.Encode()
	f := &Fetcher{
		uid:    "",
		rawURL: u,
	}
	return f, f.test()
}
