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
		if err == nil {
			port = newport
		}
	}
	log.Print("Start proxy")
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), proxy))
}
