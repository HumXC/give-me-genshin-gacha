# README

## 关于

这是一个使用 wails v2 + vue3 编写的原神祈愿记录查看工具
因为想找点事做所以就写了

仅支持 windows 系统
目前仅支持 简体中文(zh-cn) 的 大陆地区官服的原神。
如果有人用，想要过多支持或者功能可以发 issues，我看到了会尝试实现

## 截图

![截图](https://raw.githubusercontent.com/HumXC/give-me-genshin-gacha/main/doc/1.png)

## 编译

部署好 wails v2 的开发环境
在项目根目录下运行 `wails build` 编译成二进制文件
更多编译参数可以查看 wails 的官方文档:
https://wails.io/zh-Hans/docs/reference/cli#构建

# 运行

如果你是 Windows11 系统或者安装了 Edge 浏览器，一般情况下双击就能运行了。
不过，因为依赖 Webview2，如果在你的系统里找不到 Webview2，程序会让你选择是否下载 Webview2。
如果你不下载就无法运行，这也是此程序只有 10 几 mb 大小的一个原因

## 其他

我参考了以下两个项目来实现功能
https://github.com/sunfkny/genshin-gacha-export
https://github.com/biuuu/genshin-wish-export

程序内的物品图标来自提瓦特大地图
https://webstatic.mihoyo.com/ys/app/interactive-map

实现代理服务器抓取链接使用的库
https://github.com/lqqyt2423/go-mitmproxy
