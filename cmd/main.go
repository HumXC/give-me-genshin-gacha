package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

// 搜索游戏日志获取游戏数据文件的目录
func GetGameDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("无法获取用户目录: " + err.Error())
	}
	// 读取原神日志文件
	logFileName := path.Join(homeDir, "AppData", "LocalLow", "miHoYo", "原神", "output_log.txt")
	logFile, err := os.Open(logFileName)
	if err != nil {
		return "", errors.New("无法读取游戏日志: " + err.Error())
	}
	defer logFile.Close()
	logScanner := bufio.NewScanner(logFile)
	// 获取游戏数据目录名称
	searchName := "YuanShen_Data"
	for {
		if !logScanner.Scan() {
			break
		}
		line := logScanner.Text()
		err := logScanner.Err()
		if err != nil {
			return "", errors.New("日志解析错误: " + err.Error())
		}
		if !strings.Contains(line, "Warmup file") {
			continue
		}
		if !strings.Contains(line, searchName) {
			continue
		}

		i := strings.LastIndex(line, searchName)
		return line[12 : i+len(searchName)], nil
	}
	return "", errors.New("罕见错误, 没有找到游戏目录")
}

// 从游戏目录中的网络缓存获取旅行者祈愿的 URL
func GetRawURL(gameDataDir string) (string, error) {
	// 读取网络日志
	// TODO: 直接读取，而不是先使用 powershell 复制，powershell 启动缓慢
	webCacheName := path.Join(gameDataDir, "webCaches", "Cache", "Cache_Data", "data_2")
	exec.Command("powershell.exe", "/C", "Copy-Item", "\""+webCacheName+"\"", "temp").Output()

	webCache, err := os.ReadFile("temp")
	if err != nil {
		return "", errors.New("读取缓存失败: " + err.Error())
	}
	// os.Remove("temp")
	// temp 的数据由 “0” 分割
	// 提取出 temp 里的 urll 字符串
	var strEnd int

	api := "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog"
	prefx := "1/0/"
	for i := len(webCache) - 1; i > 0; i-- {
		b := webCache[i]
		if b != 0 {
			if strEnd == 0 {
				strEnd = i
			}
			continue
		}

		if strEnd == 0 {
			continue
		}

		// 将数据以 “0” 分段
		str := string(webCache[i+1 : strEnd+1])
		strEnd = 0
		// 是否为链接，链接在 temp 里以 “1/0/” 开头
		if !strings.HasPrefix(str, prefx) {
			continue
		}
		// 检查是否为祈愿记录 api 的 url
		if !strings.HasPrefix(str, prefx+api) {
			continue
		}
		return str[4:], nil

	}
	return "", errors.New("没有找到祈愿链接，尝试在游戏里打开祈愿历史记录页面")
}
func main() {
	gameDataDir, err := GetGameDir()
	if err != nil {
		log.Fatal("获取游戏目录时异常: ", err)
		return
	}
	rawURL, err := GetRawURL(gameDataDir)
	if err != nil {
		log.Fatal("解析游戏缓存时异常: ", err)
		return
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
		fmt.Printf("正在获取[%s]第 %d 页\n", ParseGachaType(gachaType), page)
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
