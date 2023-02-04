package gacha

// 提供系统代理相关的功能
import (
	"context"
	"errors"
	"net"
	"strings"
	"sync"

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
	proxy     *proxy.Proxy
	err       chan error
	cancel    context.CancelFunc
	finder    *urlFinder
	isStarted bool // 是否有正在获取 URL 的任务
	isClosed  bool
	info      systemProxyInfo
	mu        sync.Mutex
}

// 开启代理服务器并尝试捕获 URL
// targetURL 是希望匹配的 URL 列表
// 代理服务器会以 strings.HasPrefix(url, targetURL) 的方式匹配 URL
// 获取到 URL 后自动调用 Stop(), 因为代理服务器无法响应祈愿页面的 https 请求
// 感觉就像是断网了, 但是只要获取到了链接, 那么目的就达到了
// 在调用 Close() 之前, Start() 可以多次调用
// 如果在其他协程中调用了 Stop() 或 Close(), 将返回空字符串和 nil
func (p *ProxyServer) Start(targetURL string) (string, error) {
	if p.isStarted {
		return "", errors.New("代理服务器正在进行任务")
	}
	if p.isClosed {
		return "", errors.New("代理服务器已被销毁，需要重新创建")
	}
	info, err := GetSystemProxyInfo()
	if err != nil {
		return "", err
	}
	p.info = info
	p.finder.TargetURL = targetURL
	err = gosysproxy.SetGlobalProxy(p.proxy.Opts.Addr, info.Override)
	if err != nil {
		return "", err
	}
	gosysproxy.Flush()
	if err != nil {
		return "", err
	}
	p.isStarted = true
	defer p.Stop()
	p.finder.SavedURL = make(chan string)
	select {
	case err := <-p.err:
		return "", err
	case savedURL := <-p.finder.SavedURL:
		return savedURL, nil
	}
}

// 恢复系统原有的代理设置, 一般用于中断 Start() 方法
func (p *ProxyServer) Stop() error {
	if !p.isStarted {
		return nil
	}
	err := gosysproxy.SetGlobalProxy(p.info.Server, p.info.Override)
	if err != nil {
		return err
	}
	err = gosysproxy.Flush()
	if err != nil {
		return err
	}
	if !p.info.IsEnable {
		err := gosysproxy.Off()
		if err != nil {
			return err
		}
	}
	close(p.finder.SavedURL)
	p.isStarted = false
	return nil
}

// 关闭代理服务器
func (p *ProxyServer) Close() (err error) {
	// 上锁确保 isClosed 的值在各个协程中同步
	p.mu.Lock()
	if p.isClosed {
		return nil
	}
	if p.isStarted {
		err = p.Stop()
		if err != nil {
			return
		}
	}
	p.cancel()
	err = p.proxy.Close()
	if err != nil {
		p.isClosed = true
		return
	}
	p.isClosed = true
	p.mu.Unlock()
	return
}

// 创建一个代理服务器，只用于捕获请求的 URL
func NewProxyServer() (*ProxyServer, error) {
	freeAddr, err := GetFreeAddr()
	if err != nil {
		return nil, err
	}
	opts := &proxy.Options{
		Debug:             0,
		Addr:              freeAddr,
		StreamLargeBodies: 1024 * 1024 * 5,
	}
	pro, err := proxy.NewProxy(opts)
	if err != nil {
		return nil, err
	}
	addon := &urlFinder{}
	pro.AddAddon(addon)
	p := &ProxyServer{
		proxy:  pro,
		finder: addon,
		err:    make(chan error, 1),
	}
	ctx, cancel := context.WithCancel(context.Background())
	p.cancel = cancel
	go func() {
		select {
		case <-ctx.Done():
			close(p.err)
		// TODO: 测试此处是否存在协程泄露
		case p.err <- p.proxy.Start():
			if p.err != nil {
				p.cancel = nil
			}
		}
	}()
	return p, nil
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
