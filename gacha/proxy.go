package gacha

import (
	"strings"

	"github.com/Trisia/gosysproxy"
	"github.com/lqqyt2423/go-mitmproxy/proxy"
	"golang.org/x/sys/windows/registry"
)

type AddHeader struct {
	proxy.BaseAddon
	url chan (string)
}

func (a *AddHeader) Requestheaders(f *proxy.Flow) {
	url := f.Request.URL.String()
	if strings.HasPrefix(url, Api) {
		a.url <- url
		close(a.url)
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
	proxy              *proxy.Proxy
	addon              *AddHeader
	IsRunning          bool
	defaultProxyAddr   string
	defaultProxyEnable bool
}

func (p *ProxyServer) Start() error {
	addr, _ := SystemProxyAddr()
	enable, _ := IsEnableSystemProxy()
	p.defaultProxyAddr = addr
	p.defaultProxyEnable = enable

	opts := &proxy.Options{
		Debug:             0,
		Addr:              ":8080",
		StreamLargeBodies: 1024 * 1024 * 5,
	}
	pro, err := proxy.NewProxy(opts)
	if err != nil {
		return err
	}
	addon := AddHeader{}
	p.addon = &addon
	pro.AddAddon(&addon)
	p.proxy = pro
	err = p.Start()
	if err != nil {
		return err
	}
	p.IsRunning = true
	err = gosysproxy.SetGlobalProxy(addr)
	if err != nil {
		return err
	}
	gosysproxy.Flush()
	if err != nil {
		return err
	}
	return nil
}

func (p *ProxyServer) Stop() error {
	if !p.defaultProxyEnable {
		err := gosysproxy.Off()
		if err != nil {
			return err
		}
	}
	err := gosysproxy.SetGlobalProxy(p.defaultProxyAddr)
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
