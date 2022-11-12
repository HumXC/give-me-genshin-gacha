package main

import (
	"bufio"
	"errors"
	"give-me-genshin-gacha/database"
	"give-me-genshin-gacha/network"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

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
	return "", errors.New("没有找到游戏目录, 尝试进入游戏后再运行")
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

// TODO: 增加 以系统代理获取url的功能
// TODO: 后端
// TODO: 前端
func main() {
	gameDataDir, err := GetGameDir()
	if err != nil {
		log.Fatal("获取游戏目录时异常: ", err)
		return
	}
	log.Println("找到游戏目录: ", gameDataDir)
	rawURL, err := GetRawURL(gameDataDir)
	if err != nil {
		log.Fatal("解析游戏缓存时异常: ", err)
		return
	}

	fetcher, err := network.NewFetcher(rawURL)
	if err != nil {
		log.Fatal("爬虫创建失败:", err)
	}
	db, err := database.NewDB()
	if err != nil {
		log.Fatal("数据库创建失败", err)
	}

	err = SyncFetcherToDB(fetcher, db)
	if err != nil {
		log.Fatal("同步失败:", err)
	}
}
func SyncFetcherToDB(f *network.Fecher, db *database.DB) error {
	lastIDs, err := db.GetLastIDs()
	if err != nil {
		return err
	}
	for _, v := range network.GachaType {
		err := f.Get(v, lastIDs)
		if err != nil {
			return err
		}
	}
	for _, v := range f.Result {
		err := db.Add(v)
		if err != nil {
			return err
		}
	}
	return nil
}
