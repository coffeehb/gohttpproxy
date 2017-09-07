package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/google/martian"
	"golang.org/x/net/proxy"
	"syscall"
)

var (
	addr    = flag.String("addr", ":8080", "host:port of the proxy")
	forward = flag.String("forward", "", "forward to upstream proxy, example: socks5://127.0.0.1:1080")
)

func main() {
	p := martian.NewProxy()
	defer p.Close()

	fmt.Printf("Now will connect parent proxy: %v \n", *forward)
	if *forward != "" {
		url, err := url.Parse(*forward)
		if err != nil {
			fmt.Printf("forward url.Parse failed: %v\n", err)
			os.Exit(-1)
		}
		px, err := proxy.FromURL(url, proxy.Direct)
		tr := &http.Transport{
			Dial:                  px.Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		p.SetRoundTripper(tr)
	} else {

		tr := &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		p.SetRoundTripper(tr)
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting proxy on : %s", l.Addr().String())

	go p.Serve(l)

	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-sigc

	log.Println("martian: shutting down")
	os.Exit(0)
}

func init() {
	martian.Init()
}
