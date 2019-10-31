package main

import (
	"crypto/tls"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/google/martian/log"

	"syscall"

	"github.com/google/martian"
)

var (
	addr = flag.String("addr", ":8080", "host:port of the proxy")
	lv   = flag.Int("lv", log.Debug, "default log level")
)

func main() {
	p := martian.NewProxy()
	defer p.Close()
	//设置默认级别
	log.SetLevel(*lv)

	tr := &http.Transport{
		IdleConnTimeout:       300 * time.Second,
		TLSHandshakeTimeout:   3 * time.Second,
		ExpectContinueTimeout: 3 * time.Second,
		MaxIdleConns:          128,
		MaxIdleConnsPerHost:   128,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	p.SetDial((&net.Dialer{
		KeepAlive: 300 * time.Second,
	}).Dial)
	p.SetRoundTripper(tr)

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Errorf(err.Error())
	}

	log.Infof("Starting proxy on : %s", l.Addr().String())

	go p.Serve(l)

	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-sigc

	log.Infof("Notice: shutting down")
	os.Exit(0)
}

func init() {
	martian.Init()
}
