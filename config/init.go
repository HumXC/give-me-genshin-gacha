// 提供配置文件的序列/反序列和存储过程
package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type savedURL struct {
	Uid string `json:"uid"`
	URL string `json:"url"`
}

type Config struct {
	filePath    string
	SavedURLs   []savedURL `json:"savedUrls"` // 祈愿页面的 URL
	Language    string     `json:"language"`
	SelectedUid string     `json:"selectedUid"`
	ShowGacha   []string   `json:"showGacha"`
	IsDarkTheme bool       `json:"isDarkTheme"`
	IsAutoSync  bool       `json:"isAutoSync"`
	IsUseProxy  bool       `json:"isUseProxy"`
}

var config *Config

func GetConfig(filePath string) (*Config, error) {
	if config != nil {
		return config, nil
	}
	c := Config{
		filePath:  filePath,
		SavedURLs: make([]savedURL, 0),
		ShowGacha: make([]string, 0),
	}
	if f, err := os.Open(filePath); err == nil {
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		conf := Config{}
		err = json.Unmarshal(b, &conf)
		if err != nil {
			return nil, err
		}
		conf.filePath = c.filePath
		c = conf
	}
	config = &c
	return &c, nil
}

// 将配置写入文件
func (c *Config) Save() error {
	b, err := json.Marshal(*c)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, "", "    ")
	if err != nil {
		return err
	}
	f, err := os.OpenFile(c.filePath, os.O_CREATE|os.O_WRONLY, 0775)
	if err != nil {
		return err
	}
	defer f.Close()
	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return err
}

// 提供一个 Config cfg, 把 cfg 的变量值同步到自身
func (c *Config) Put(cfg Config) {

}
