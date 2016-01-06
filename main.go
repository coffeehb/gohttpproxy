package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Print("Start proxy")
	to := 3 * time.Second
	proxy := goproxy.NewProxyHttpServer()
	proxy.Tr.Dial = func(network, addr string) (c net.Conn, err error) {
		c, err = (&net.Dialer{
			Timeout:   to,
			KeepAlive: to,
		}).Dial(network, addr)
		if c, ok := c.(*net.TCPConn); err == nil && ok {
			log.Println("Set keep alive")
			c.SetKeepAlive(true)
			c.SetNoDelay(true)
			c.SetDeadline(time.Now().Add(to))
		}
		return
	}
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe("127.0.0.1:8123", proxy))
}
