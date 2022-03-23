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

	//设置默认级别
	log.SetLevel(*lv)

	p := martian.NewProxy()
	//设置读写超时为600分钟，也就是10小时
	p.SetTimeout(600 * time.Minute)
	defer p.Close()

	tr := &http.Transport{
		IdleConnTimeout:       75 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 5 * time.Second,
		DisableKeepAlives:     true,
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   0,
		MaxConnsPerHost:       4096,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	p.SetDial((&net.Dialer{
		KeepAlive: -1,
		Timeout:   5 * time.Second,
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
