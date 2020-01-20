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
)

func main() {
	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()
	p := martian.NewProxy()
	defer p.Close()
	//设置默认级别
	log.SetLevel(*lv)

	tr := &http.Transport{
		IdleConnTimeout:       60 * time.Second,
		ResponseHeaderTimeout: 4 * time.Second,
		TLSHandshakeTimeout:   4 * time.Second,
		ExpectContinueTimeout: 4 * time.Second,
		MaxIdleConns:          128,
		MaxIdleConnsPerHost:   128,
		MaxConnsPerHost: 128,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	p.SetDial((&net.Dialer{
		KeepAlive: 60 * time.Second,
		Timeout: 4 * time.Second,
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
