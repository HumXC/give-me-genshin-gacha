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

var gachaTypes = []string{
	GachaType100, GachaType200, GachaType301, GachaType302,
}

func GachaTypes() []string {
	result := make([]string, len(gachaTypes))
	copy(result, gachaTypes)
	return result
}

var ErrURLTestFailed = errors.New("url test failed")
var ErrURLAuthTimeout = errors.New("url authkey timeout")
var ErrPageEnd = errors.New("page end")

type RawGachaLog struct {
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

type GachaLogFetcher struct {
	uid       uint64
	lang      string
	gachaType string
	lock      sync.Mutex
	rawURL    *url.URL
	haveNext  bool
	targetID  string
	endID     string
}

// 获取此祈愿链接的 Uid
func (f *GachaLogFetcher) Uid() uint64 {
	return f.uid
}

// 获取此祈愿链接的 Uid
func (f *GachaLogFetcher) Lang() string {
	return f.lang
}

func (f *GachaLogFetcher) SetGachaType(gachaType, endID string) {
	f.gachaType = gachaType
	f.targetID = endID
	f.endID = "0"
	f.haveNext = true
}

func (f *GachaLogFetcher) Next() bool {
	return f.haveNext
}

func (f *GachaLogFetcher) Logs() ([]RawGachaLog, error) {
	if !f.haveNext {
		return []RawGachaLog{}, nil
	}
	result, err := f.fetch(f.gachaType, f.endID)
	if err != nil {
		return result, err
	}
	length := len(result)
	if length == 0 {
		f.haveNext = false
		return result, nil
	}
	if length < 20 {
		f.haveNext = false
	}
	f.endID = result[length-1].ID
	// 筛选出更新的物品
	for i, item := range result {
		if item.ID == f.targetID {
			f.haveNext = false
			return result[:i], ErrPageEnd
		}
	}
	return result, nil
}

// 获取指定祈愿的所有记录，gachaType 是数字代号的字符串
func (f *GachaLogFetcher) fetch(gachaType, endID string) ([]RawGachaLog, error) {
	type ResponseData struct {
		Page  string        `json:"page"`
		Size  string        `json:"size"`
		Total string        `josn:"total"`
		List  []RawGachaLog `json:"list"`
	}
	type Response struct {
		RetCode int          `json:"retcode"`
		Message string       `json:"message"`
		Data    ResponseData `json:"data"`
		Region  string       `json:"region"`
	}
	result := make([]RawGachaLog, 0)
	f.lock.Lock()
	defer f.lock.Unlock()
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
	return result, nil
}

func (f *GachaLogFetcher) test() error {
	gachaTypes := GachaTypes()
	for i := 0; i < len(gachaTypes); i++ {
		logs, err := f.fetch(gachaTypes[i], "0")
		if err != nil {
			return fmt.Errorf("URL 测试失败: %w", err)
		}
		if len(logs) == 0 {
			if i == len(gachaTypes)-1 {
				return errors.New("URL 没有祈愿记录")
			}
			continue
		}
		uid, _ := strconv.ParseUint(logs[0].Uid, 10, 64)
		f.uid = uid
		break
	}
	return nil
}

// 根据祈愿的链接创建爬虫
func NewGachaLogFetcher(rawURL string, lang string) (*GachaLogFetcher, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("size", "20")
	query.Set("page", "1")
	query.Set("begin_id", "0")
	query.Set("lang", lang)
	u.RawQuery = query.Encode()
	f := &GachaLogFetcher{
		rawURL: u,
	}
	return f, f.test()
}
