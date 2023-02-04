package assets

import (
	"net/http"
	"os"
	"path"
	"strings"
)

// 提供动态资产的能力

// 通过 Server 保存的文件均在此路径下
const CHAHE_DIR string = "./cache"

// 资源类型:
// 祈愿物品的图标 /icon/gacha_item/:id
type Server struct {
	iconHandler *iconHandler
}

func NewServer() (*Server, error) {
	f := Server{
		iconHandler: &iconHandler{
			IconDir: path.Join(CHAHE_DIR, "icons"),
		},
	}
	if !IsExist(f.iconHandler.IconDir) {
		err := os.MkdirAll(f.iconHandler.IconDir, 0755)
		if err != nil {
			return nil, err
		}
	}
	return &f, nil
}

func (f *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	isItemIconReq := strings.HasPrefix(req.URL.Path, "/icon/gacha_item/")
	if isItemIconReq {
		f.iconHandler.ServeHTTP(resp, req)
		return
	}
}
