package main

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"os"
	"strconv"
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
	log.SetOutput(os.Stdout)
	log.Print(http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), proxy))
}
