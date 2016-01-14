package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	var host string
	var port int
	args := os.Args[1:]
	argslen := len(args)
	fmt.Println(argslen)
	host = "127.0.0.1"
	port = 8123
	if argslen >= 1 {
		host = args[0]
	}
	if argslen >= 2 {
		newport, err := strconv.Atoi(args[1])
		if err != nil {
			port = 53
		} else {
			port = newport
		}
	}
	log.Print("Start proxy")
	proxy := goproxy.NewProxyHttpServer()
	to := 3 * time.Second
	proxy.Tr.Dial = func(network, addr string) (c net.Conn, err error) {
		c, err = (&net.Dialer{
			Timeout:   to,
			KeepAlive: to,
		}).Dial(network, addr)
		if c, ok := c.(*net.TCPConn); err == nil && ok {
			log.Println("Set keep alive")
			c.SetKeepAlive(true)
			c.SetNoDelay(true)
		}
		return
	}
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), proxy))
}
