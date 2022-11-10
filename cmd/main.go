package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	角色活动祈愿 = "301"
	武器活动祈愿 = "302"
	常驻祈愿   = "200"
	新手祈愿   = "100"
)

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

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("无法获取用户目录")
		return
	}
	// 读取原神日志文件
	logFileName := path.Join(homeDir, "AppData", "LocalLow", "miHoYo", "原神", "output_log.txt")
	logFile, err := os.Open(logFileName)
	if err != nil {
		fmt.Println("无法读取游戏日志")
		return
	}
	defer logFile.Close()
	logScanner := bufio.NewScanner(logFile)
	// 获取游戏数据目录名称
	gameDataDir := "YuanShen_Data"
	for {
		logScanner.Scan()
		line := logScanner.Text()
		err := logScanner.Err()
		if err != nil {
			fmt.Println("日志解析错误, 尝试进入游戏")
			return
		}
		if !strings.Contains(line, "Warmup file") {
			continue
		}
		if !strings.Contains(line, gameDataDir) {
			continue
		}

		i := strings.LastIndex(line, gameDataDir)
		gameDataDir = line[12 : i+len(gameDataDir)]
		fmt.Println("找到游戏数据目录: ", gameDataDir)
		break
	}

	// 读取网络日志
	// TODO: 直接读取，而不是先使用 powershell 复制，powershell 启动缓慢
	webCacheName := path.Join(gameDataDir, "webCaches", "Cache", "Cache_Data", "data_2")
	exec.Command("powershell.exe", "/C", "Copy-Item", "\""+webCacheName+"\"", "temp").Output()

	webCache, err := os.ReadFile("temp")
	if err != nil {
		fmt.Println("读取缓存失败: ", err)
		return
	}
	// os.Remove("temp")
	// temp 的数据由 “0” 分割
	// 提取出 temp 里的 urll 字符串
	var rawURL string
	var urlEnd int

	api := "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog"

	for i := len(webCache) - 1; i > 0; i-- {
		b := webCache[i]
		if b != 0 {
			if urlEnd == 0 {
				urlEnd = i
			}
			continue
		}

		if urlEnd == 0 {
			continue
		}

		str := string(webCache[i+1 : urlEnd+1])
		urlEnd = 0
		// 链接在 temp 里以 “1/0/” 开头
		prefx := "1/0/"
		if !strings.HasPrefix(str, prefx) {
			continue
		}
		s := strings.TrimPrefix(str, prefx)
		if strings.HasPrefix(s, api) {
			rawURL = s
			break
		}

	}

	// 解析链接参数
	f := NewFetcher(rawURL)
	f.Get(常驻祈愿)

	fmt.Printf("%v", f)
}

type Fecher struct {
	Uid    string
	url    *url.URL
	Result map[string][]RespDataListItem
}

// 获取指定祈愿的所有记录，gachaType 是数字代号的字符串
func (f *Fecher) Get(gachaType string) error {
	list := f.Result[gachaType]
	page := 1
	endID := "0"
	if len(list) > 0 {
		endID = list[len(list)-1].ID
	}
	query := f.url.Query()

	query.Set("gacha_type", gachaType)
	for {
		fmt.Printf("正在为[%s]获取[%s]第 %d 页\n", f.url, ParseGachaType(gachaType), page)
		query.Set("page", strconv.Itoa(page))
		query.Set("end_id", endID)
		f.url.RawQuery = query.Encode()
		url_ := f.url.String()
		resp, err := http.DefaultClient.Get(url_)
		if err != nil {
			return err
		}
		jsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		var r Response
		err = json.Unmarshal(jsonData, &r)
		if err != nil {
			return err
		}
		if r.Message != "OK" {
			return errors.New("Api error: " + r.Message)
		}
		if len(r.Data.List) == 0 {
			break
		}
		page++
		endID = r.Data.List[len(r.Data.List)-1].ID
		list = append(list, r.Data.List...)
		time.Sleep(1 * time.Second)
	}
	f.Result[gachaType] = list
	return nil
}

func NewFetcher(rawURL string) *Fecher {
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("URL 解析失败: ", err)
		return nil
	}
	query := u.Query()
	query.Set("size", "20")
	query.Set("end_id", "0")
	u.RawQuery = query.Encode()
	return &Fecher{
		Uid:    query.Get("uid"),
		url:    u,
		Result: make(map[string][]RespDataListItem),
	}

}
