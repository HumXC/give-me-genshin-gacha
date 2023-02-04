package assets

import (
	"errors"
	"fmt"
	"give-me-genshin-gacha/gacha"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

type Icon struct {
	Url  string `json:"icon"`
	Name string `json:"name"`
}

// 为 Server 提供处理 icon 请求的能力
type iconHandler struct {
	IconDir     string
	assetsStore gacha.AssetsStore
	isFetching  sync.Mutex
}

func (i *iconHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	icon := strings.TrimPrefix(req.URL.Path, "/icon/gacha_item/")
	fileName := path.Join(i.IconDir, icon)
	if !IsExist(fileName) {
		http.Error(resp, "你先别急，后台正在获取图片", http.StatusNotFound)
		if i.isFetching.TryLock() {
			go i.fetch()
		}
		return
	}
	b, err := os.ReadFile(fileName)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = resp.Write(b)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (i *iconHandler) fetch() {
	defer func() {
		i.isFetching.Unlock()
	}()
	fmt.Println("正在获取图标资源")
	items, err := i.assetsStore.Get("zh-cn")
	if err != nil {
		fmt.Println("获取图标资源失败: " + err.Error())
		return
	}
	var wg sync.WaitGroup
	ch := make(chan struct{}, 16)
	for j := 0; j < len(items); j++ {
		fileName := path.Join(i.IconDir, strconv.Itoa(items[j].ItemID)+".png")
		if IsExist(fileName) {
			continue
		}
		ch <- struct{}{}
		wg.Add(1)
		item := items[j]
		go func() {
			defer wg.Done()
			saveWebTo(item.Icon, fileName)
			<-ch
		}()
	}
	wg.Wait()
	fmt.Println("获取图标资源完毕")
}

func saveWebTo(url, fileName string) error {
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	f, err := os.Create(fileName)
	if err != nil {
		_ = f.Close()
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		_ = f.Close()
		_ = os.Remove(f.Name())
		return err
	}
	_ = f.Close()
	return nil
}

// 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || errors.Is(err, os.ErrExist)
}
