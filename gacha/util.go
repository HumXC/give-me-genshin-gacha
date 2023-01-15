package gacha

import (
	"bufio"
	"fmt"
	"give-me-genshin-gacha/models"
	"os"
	"path"
	"strings"
	"syscall"
)

const Api = "https://hk4e-api.mihoyo.com/event/gacha_info/api/getGachaLog"

// 搜索游戏日志获取游戏数据文件的目录
func GetGameDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户失败目录: %w", err)
	}
	// 读取原神日志文件
	logFileName := path.Join(homeDir, "AppData", "LocalLow", "miHoYo", "原神", "output_log.txt")
	logFile, err := os.Open(logFileName)
	if err != nil {
		return "", fmt.Errorf("读取游戏日志失败: %w", err)
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
			return "", fmt.Errorf("日志解析错误: %w", err)
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
	return "", fmt.Errorf("没有找到游戏目录, 进入游戏后再尝试")
}

// 从游戏目录中的网络缓存获取旅行者祈愿的 URL
func GetRawURL(gameDataDir string) (string, error) {
	// 读取网络日志
	webCacheP, err := GetWebCacha(gameDataDir)
	if err != nil {
		return "", fmt.Errorf("读取缓存失败: %w", err)
	}
	webCache := *webCacheP
	// temp 的数据由 “0” 分割
	// 提取出 temp 里的 urll 字符串
	var strEnd int
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
		if !strings.HasPrefix(str, prefx+Api) {
			continue
		}
		return str[4:], nil

	}
	return "", fmt.Errorf("没有找到祈愿链接，尝试在游戏里打开祈愿历史记录页面")
}

// 读取游戏目录内的网络缓存
func GetWebCacha(gameDataDir string) (*[]byte, error) {
	//参考，搜了很久。tnnd
	// https://github.com/golang/go/issues/46164
	// http://zplutor.github.io/2018/08/26/file-share-mode-and-access-rights/
	result := make([]byte, 2048)
	fileName := path.Join(gameDataDir, "webCaches", "Cache", "Cache_Data", "data_2")
	ptr, err := syscall.UTF16PtrFromString(fileName)
	if err != nil {
		return nil, err
	}
	f, err := syscall.CreateFile(ptr, syscall.GENERIC_READ|syscall.GENERIC_WRITE, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(f)
	b := make([]byte, 1024)
	for {
		n, err := syscall.Read(f, b)
		if n == 0 {
			return &result, nil
		}
		if err != nil {
			return nil, err
		}
		result = append(result, b[:n]...)
	}
}

// 将爬虫的数据转换成数据库的数据
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
