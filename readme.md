# gohttpproxy

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
