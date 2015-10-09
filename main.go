package main

import (
	"flag"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func main() {
	verbose := true
	addr := "127.0.0.1:8123"
	flag.Parse()
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose
	log.Fatal(http.ListenAndServe(addr, proxy))
}
