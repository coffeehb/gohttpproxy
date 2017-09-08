# gohttpproxy
Update 2017-09-07

now gohttpproxy support upstream proxy, like socks5

```
browser => gohttpproxy => socks5
```

Go http(s) proxy , By default listen on 127.0.0.1:8123


```
Usage of ./gohttpproxy:
  -addr string
        host:port of the proxy (default ":8080")
  -forward string
        forward to upstream proxy, example: socks5://127.0.0.1:1080
  -v int
        log level

```

## Install


```
go get -u -v github.com/golang/dep/cmd/dep
dep ensure -v
go build -buildmode=pie -v
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
