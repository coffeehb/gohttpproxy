package main

import (
	"crypto/tls"
	"flag"
	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "host:port of the proxy")
	lv   = flag.Int("lv", log.Debug, "default log level")
	h    = flag.Bool("h", false, "help")
)

func main() {
	flag.Parse()
	if *h {
		flag.PrintDefaults()
		os.Exit(0)
	}

	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()
	p := martian.NewProxy()
	//设置读写超时为600分钟，也就是10小时
	p.SetTimeout(600 * time.Minute)
	defer p.Close()
	//设置默认级别
	log.SetLevel(*lv)

	tr := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		ResponseHeaderTimeout: 3 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 0,
		DisableKeepAlives: true,
		MaxIdleConns:          -1,
		MaxIdleConnsPerHost:   -1,
		MaxConnsPerHost:       512,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	p.SetDial((&net.Dialer{
		KeepAlive: 75 * time.Second,
		Timeout:   3 * time.Second,
	}).Dial)
	p.SetRoundTripper(tr)

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Errorf(err.Error())
	}

	log.Infof("Starting proxy on : %s", l.Addr().String())

	go p.Serve(l)

	signChannel := make(chan os.Signal, 2)
	signal.Notify(signChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-signChannel

	log.Infof("Notice: shutting down")
	os.Exit(0)
}

func init() {
	martian.Init()
}
