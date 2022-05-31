# gohttpproxy

一款小巧，但强大的代理服务器软件，单文件部署，占用资源少，使用简便。吃得少，干的活多。稳定，坚强，历经7年多时间持续打磨。一直在迭代一直在更新。
如今更好用。

可以从release页面下载最新的已经编译好的二进制文件部署，支持windows, mac, linux全部平台。

https://github.com/cnmade/gohttpproxy/releases/tag/v2.0.0


go http proxy is a simple http proxy, support `HTTP CONNECT ` proxy


```
browser => gohttpproxy => target web site
```

Go http(s) proxy , By default listen on 127.0.0.1:8123


```
Usage of ./gohttpproxy:
  -addr string
        host:port of the proxy (default ":8080")
  -lv int
        log level: 1: debug, 2: info, 3: debug

```

## Install


``` 
CGO_ENABLED=0 go build -v -a -ldflags ' -s -w  -extldflags "-static"' .
# go1.14rc1
CGO_ENABLED=0 go1.14rc1 build -v -a -ldflags ' -s -w  -extldflags "-static"' .

./gohttpproxy
```
## Donate me please

### Bitcoin donate

```
136MYemy5QmmBPLBLr1GHZfkES7CsoG4Qh
```
### Alipay donate
![Scan QRCode donate me via Alipay](https://www.netroby.com/assets/images/alipayme.jpg)

**Scan QRCode donate me via Alipay**
