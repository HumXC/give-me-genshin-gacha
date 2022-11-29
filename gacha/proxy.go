package gacha

import (
	"net"
	"strings"

	"github.com/Trisia/gosysproxy"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"golang.org/x/sys/windows/registry"
)

type addHeader struct {
	proxy.BaseAddon
	url chan (string)
}

func (a *addHeader) Requestheaders(f *proxy.Flow) {
	url := f.Request.URL.String()
	if strings.HasPrefix(url, Api) {
		a.url <- url
	}
}
func openProxyReg() (registry.Key, error) {
	return registry.OpenKey(registry.CURRENT_USER,
		"Software\\Microsoft\\Windows\\CurrentVersion\\Internet Settings", registry.ALL_ACCESS)
}
func SystemProxyAddr() (string, error) {
	key, err := openProxyReg()
	if err != nil {
		return "", nil
	}
	v, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		return "", nil
	}
	return v, nil
}
func SystemProxyOverride() (string, error) {
	key, err := openProxyReg()
	if err != nil {
		return "", nil
	}
	v, _, err := key.GetStringValue("ProxyOverride")
	if err != nil {
		return "", nil
	}
	return v, nil
}
func IsEnableSystemProxy() (bool, error) {
	key, err := openProxyReg()
	if err != nil {
		return false, nil
	}
	v, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		return false, nil
	}
	return v != 0, nil
}

type ProxyServer struct {
	proxy                *proxy.Proxy
	Url                  chan (string)
	IsRunning            bool
	defaultProxyAddr     string
	defaultProxyEnable   bool
	defaultProxyOverride string
}

func (p *ProxyServer) Start() error {
	addr, _ := SystemProxyAddr()
	enable, _ := IsEnableSystemProxy()
	override, _ := SystemProxyOverride()
	p.defaultProxyAddr = addr
	p.defaultProxyEnable = enable
	p.defaultProxyOverride = override
	go func() {
		// TODO: 错误处理
		p.proxy.Start()
	}()
	err := gosysproxy.SetGlobalProxy(p.proxy.Opts.Addr, override)
	if err != nil {
		return err
	}
	gosysproxy.Flush()
	if err != nil {
		return err
	}
	p.IsRunning = true
	return nil
}

func (p *ProxyServer) Stop() error {
	defer recover()
	close(p.Url)
	if !p.defaultProxyEnable {
		err := gosysproxy.Off()
		if err != nil {
			return err
		}
	}
	err := gosysproxy.SetGlobalProxy(p.defaultProxyAddr, p.defaultProxyOverride)
	if err != nil {
		return err
	}
	err = gosysproxy.Flush()
	if err != nil {
		return err
	}
	err = p.proxy.Close()
	if err != nil {
		return err
	}
	p.IsRunning = false
	p.proxy = nil
	return nil
}

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
	addon := &addHeader{
		url: make(chan string),
	}
	pro.AddAddon(addon)

	p := &ProxyServer{
		Url:   addon.url,
		proxy: pro,
	}
	return p, nil
}
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
