package gacha

// 提供系统代理相关的功能
import (
	"context"
	"net"
	"strings"

	"github.com/Trisia/gosysproxy"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"golang.org/x/sys/windows/registry"
)

// go-mitmproxy/proxy 的插件
// 查找指定的 URL 然后返回
type urlFinder struct {
	proxy.BaseAddon
	SavedURL  chan string
	TargetURL string
	Ctx       context.Context
}

func (a *urlFinder) Requestheaders(f *proxy.Flow) {
	url := f.Request.URL.String()
	if strings.HasPrefix(url, a.TargetURL) {
		a.SavedURL <- url
	}
}

type systemProxyInfo struct {
	Server   string
	IsEnable bool
	Override string
}

// 打开注册表获取系统代理的相关信息
func openProxyReg() (registry.Key, error) {
	return registry.OpenKey(registry.CURRENT_USER,
		"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", registry.ALL_ACCESS)
}

// 获取当前系统代理的信息
func GetSystemProxyInfo() (result systemProxyInfo, e error) {
	key, err := openProxyReg()
	if err != nil {
		return
	}
	server, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		e = err
		return
	}
	override, _, err := key.GetStringValue("ProxyOverride")
	if err != nil {
		e = err
		return
	}
	isEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		e = err
		return
	}
	result.Server = server
	result.Override = override
	result.IsEnable = isEnable != 0
	return result, nil
}

type ProxyServer struct {
	savedURL chan string
	err      chan error
	Err      <-chan error
}

func (p *ProxyServer) Url() string {
	return <-p.savedURL
}

// 创建一个代理服务器，只用于捕获请求的 URL
func NewProxyServer(ctx context.Context, targetURL string) *ProxyServer {
	urlChan := make(chan string, 1)
	errChan := make(chan error, 1)
	server := ProxyServer{
		savedURL: urlChan,
		err:      errChan,
		Err:      errChan,
	}
	freeAddr, err := GetFreeAddr()
	if err != nil {
		close(urlChan)
		errChan <- err
		return &server
	}
	opts := &proxy.Options{
		Debug:             0,
		Addr:              freeAddr,
		StreamLargeBodies: 1024 * 1024 * 5,
	}
	pro, err := proxy.NewProxy(opts)
	if err != nil {
		close(urlChan)
		server.err <- err
		return &server
	}

	addon := &urlFinder{
		SavedURL:  urlChan,
		TargetURL: targetURL,
	}
	pro.AddAddon(addon)

	proErr := make(chan error)
	go func() {
		err := pro.Start()
		if err != nil {
			proErr <- err
		}
		close(proErr)
	}()
	go func() {
		// 系统默认的代理配置
		info, err := GetSystemProxyInfo()
		if err != nil {
			errChan <- err
			return
		}
		defer close(errChan)
		defer func() {
			// 关闭系统代理
			err := gosysproxy.SetGlobalProxy(info.Server, info.Override)
			if err != nil {
				errChan <- err
				return
			}
			err = gosysproxy.Flush()
			if err != nil {
				errChan <- err
				return
			}
			if info.IsEnable {
				return
			}
			err = gosysproxy.Off()
			if err != nil {
				errChan <- err
				return
			}
		}()
		defer close(urlChan)

		// 设置系统代理
		err = gosysproxy.SetGlobalProxy(opts.Addr, info.Override)
		if err != nil {
			errChan <- err
			return
		}
		err = gosysproxy.Flush()
		if err != nil {
			errChan <- err
			return
		}
		select {
		case <-ctx.Done():
			err := pro.Close()
			if err != nil {
				errChan <- err
			}
		case err := <-proErr:
			if err != nil {
				errChan <- err
			}
		}
	}()
	return &server
}

// 获取有可用端口的 localhost 地址
func GetFreeAddr() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).String(), nil
}
