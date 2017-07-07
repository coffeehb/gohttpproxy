package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/martian"
)

var (
	addr           = flag.String("addr", ":8080", "host:port of the proxy")
)

func main() {
	p := martian.NewProxy()
	defer p.Close()

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


	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting proxy on : %s", l.Addr().String())

	go p.Serve(l)


	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	<-sigc

	log.Println("martian: shutting down")
}

func init() {
	martian.Init()
}


