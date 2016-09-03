package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/getlantern/http-proxy/commonfilter"
	"github.com/getlantern/http-proxy/filters"
	"github.com/getlantern/http-proxy/forward"
	"github.com/getlantern/http-proxy/httpconnect"
	"github.com/getlantern/http-proxy/listeners"
	"github.com/getlantern/http-proxy/server"
)

var (
	testingLocal = false

	help      = flag.Bool("help", false, "Get usage help")
	keyfile   = flag.String("key", "", "Private key file name")
	certfile  = flag.String("cert", "", "Certificate file name")
	https     = flag.Bool("https", false, "Use TLS for client to proxy communication")
	addr      = flag.String("addr", ":8080", "Address to listen")
	maxConns  = flag.Uint64("maxconns", 0, "Max number of simultaneous connections allowed connections")
	idleClose = flag.Uint64("idleclose", 30, "Time in seconds that an idle connection will be allowed before closing it")
)

func main() {
	var err error

	_ = flag.CommandLine.Parse(os.Args[1:])
	if *help {
		flag.Usage()
		return
	}

	// Logging
	// TODO: use real parameters

	filterChain := filters.Join(
		commonfilter.New(&commonfilter.Options{
			AllowLocalhost: testingLocal,
		}),
		httpconnect.New(&httpconnect.Options{IdleTimeout: time.Duration(*idleClose) * time.Second}),
		forward.New(&forward.Options{IdleTimeout: time.Duration(*idleClose) * time.Second}),
	)

	// Create server
	srv := server.NewServer(filterChain)

	// Add net.Listener wrappers for inbound connections
	srv.AddListenerWrappers(
		// Limit max number of simultaneous connections
		func(ls net.Listener) net.Listener {
			return listeners.NewLimitedListener(ls, *maxConns)
		},
		// Close connections after 30 seconds of no activity
		func(ls net.Listener) net.Listener {
			return listeners.NewIdleConnListener(ls, time.Duration(*idleClose)*time.Second)
		},
	)

	// Serve HTTP/S
	if *https {
		err = srv.ListenAndServeHTTPS(*addr, *keyfile, *certfile, nil)
	} else {
		err = srv.ListenAndServeHTTP(*addr, nil)
	}
	if err != nil {
		fmt.Println("Error serving: %v", err)
	}
}
