package assets

import (
	"context"
	"net/http"
	"os"
)

// 提供动态资产的能力

// 通过 Server 保存的文件均在此路径下
const CHAHE_DIR string = "./cache"

// 资源类型:
// 祈愿物品的图标 /icon/gacha_item/:id
type Server struct {
	ctx context.Context
}

func NewServer(ctx context.Context) *Server {
	f := Server{
		ctx: ctx,
	}

	return &f
}

func (f *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	// 分路由
}

// 判断文件是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
