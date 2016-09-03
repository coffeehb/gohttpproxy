package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/google/martian"
	"github.com/google/martian/cors"
	"github.com/google/martian/har"
	"github.com/google/martian/httpspec"
	"github.com/google/martian/marbl"
	"github.com/google/martian/martianhttp"
	"github.com/google/martian/martianlog"
	"github.com/google/martian/mitm"
	"github.com/google/martian/trafficshape"
	"github.com/google/martian/verify"

	_ "github.com/google/martian/body"
	_ "github.com/google/martian/cookie"
	mlog "github.com/google/martian/log"
	_ "github.com/google/martian/martianurl"
	_ "github.com/google/martian/method"
	_ "github.com/google/martian/pingback"
	_ "github.com/google/martian/priority"
	_ "github.com/google/martian/querystring"
	_ "github.com/google/martian/status"
)

var (
	level          = flag.Int("v", 0, "log level")
	addr           = flag.String("addr", ":8080", "host:port of the proxy")
	tlsAddr        = flag.String("tls-addr", ":4443", "host:port of the proxy over TLS")
	api            = flag.String("api", "martian.proxy", "hostname for the API")
	generateCA     = flag.Bool("generate-ca-cert", false, "generate CA certificate and private key for MITM")
	cert           = flag.String("cert", "", "CA certificate used to sign MITM certificates")
	key            = flag.String("key", "", "private key of the CA used to sign MITM certificates")
	organization   = flag.String("organization", "Martian Proxy", "organization name for MITM certificates")
	validity       = flag.Duration("validity", time.Hour, "window of time that MITM certificates are valid")
	allowCORS      = flag.Bool("cors", false, "allow CORS requests to configure the proxy")
	harLogging     = flag.Bool("har", false, "enable HAR logging API")
	trafficShaping = flag.Bool("traffic-shaping", false, "enable traffic shaping API")
	skipTLSVerify  = flag.Bool("skip-tls-verify", false, "skip TLS server verification; insecure")
)

func main() {
	flag.Parse()

	mlog.SetLevel(*level)

	p := martian.NewProxy()

	// Respond with 404 to any unknown proxy path.
	http.HandleFunc(*api+"/", http.NotFound)

	var x509c *x509.Certificate
	var priv interface{}

	if *generateCA {
		var err error
		x509c, priv, err = mitm.NewAuthority("martian.proxy", "Martian Authority", 30*24*time.Hour)
		if err != nil {
			log.Fatal(err)
		}
	} else if *cert != "" && *key != "" {
		tlsc, err := tls.LoadX509KeyPair(*cert, *key)
		if err != nil {
			log.Fatal(err)
		}
		priv = tlsc.PrivateKey

		x509c, err = x509.ParseCertificate(tlsc.Certificate[0])
		if err != nil {
			log.Fatal(err)
		}
	}

	if x509c != nil && priv != nil {
		mc, err := mitm.NewConfig(x509c, priv)
		if err != nil {
			log.Fatal(err)
		}

		mc.SetValidity(*validity)
		mc.SetOrganization(*organization)
		mc.SkipTLSVerify(*skipTLSVerify)

		p.SetMITM(mc)

		// Expose certificate authority.
		ah := martianhttp.NewAuthorityHandler(x509c)
		configure("/authority.cer", ah)

		// Start TLS listener for transparent MITM.
		tl, err := net.Listen("tcp", *tlsAddr)
		if err != nil {
			log.Fatal(err)
		}

		go p.Serve(tls.NewListener(tl, mc.TLS()))
	}

	stack, fg := httpspec.NewStack("martian")
	p.SetRequestModifier(stack)
	p.SetResponseModifier(stack)

	m := martianhttp.NewModifier()
	fg.AddRequestModifier(m)
	fg.AddResponseModifier(m)

	if *harLogging {
		hl := har.NewLogger()
		stack.AddRequestModifier(hl)
		stack.AddResponseModifier(hl)

		configure("/logs", har.NewExportHandler(hl))
		configure("/logs/reset", har.NewResetHandler(hl))
	}

	logger := martianlog.NewLogger()
	logger.SetDecode(true)

	stack.AddRequestModifier(logger)
	stack.AddResponseModifier(logger)

	lsh := marbl.NewHandler()
	http.Handle("/binlogs", lsh)

	lsm := marbl.NewModifier(lsh)
	stack.AddRequestModifier(lsm)
	stack.AddResponseModifier(lsm)

	// Proxy specific handlers.
	// These handlers take precedence over proxy traffic and will not be
	// intercepted.

	// Configure modifiers.
	configure("/configure", m)

	// Verify assertions.
	vh := verify.NewHandler()
	vh.SetRequestVerifier(m)
	vh.SetResponseVerifier(m)
	configure("/verify", vh)

	// Reset verifications.
	rh := verify.NewResetHandler()
	rh.SetRequestVerifier(m)
	rh.SetResponseVerifier(m)
	configure("/verify/reset", rh)

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}

	if *trafficShaping {
		tsl := trafficshape.NewListener(l)
		tsh := trafficshape.NewHandler(tsl)
		configure("/shape-traffic", tsh)

		l = tsl
	}

	log.Println("martian: proxy started on:", l.Addr())

	go p.Serve(l)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill)

	<-sigc

	log.Println("martian: shutting down")
}

// configure installs a configuration handler at path.
func configure(path string, handler http.Handler) {
	if *allowCORS {
		handler = cors.NewHandler(handler)
	}

	http.Handle(filepath.Join(*api, path), handler)

}
