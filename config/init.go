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
type ShowGacha struct {
	G301 bool `json:"g301"`
	G302 bool `json:"g302"`
	G200 bool `json:"g200"`
	G100 bool `json:"g100"`
}
type Config struct {
	filePath    string
	Language    string    `json:"language"`
	SelectedUid uint      `json:"selectedUid"`
	ShowGacha   ShowGacha `json:"showGacha"`
	GameDir     string    `json:"gameDir"`
	IsDarkTheme bool      `json:"isDarkTheme"`
	IsAutoSync  bool      `json:"isAutoSync"`
	IsUseProxy  bool      `json:"isUseProxy"`
}

var config *Config

func Get(filePath string) (*Config, error) {
	if config != nil {
		return config, nil
	}
	c := Config{
		filePath: filePath,
		ShowGacha: ShowGacha{
			G301: true,
			G302: true,
		},
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
	cfg.filePath = c.filePath
	// TODO: 向 wails 提交 issue
	// *c = cfg
	// 不能如此简单地赋值
	// cfg 来自前端
	// cfg 里有空切片时，此空切片将指向一个 nil，而不是一个空数组
	// 在保存 json 时，此空切片会变成 null
	// if len(cfg.SavedURLs) == 0 {
	// 	cfg.SavedURLs = make([]savedURL, 0)
	// }
	*c = cfg
}
