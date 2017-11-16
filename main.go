package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

const (
	defaultRelaySrv   = "8.8.8.8:53"
	defaultListenAddr = ":53"
)

var (
	netflixNames = []string{"netflix.com", "netflix.net", "nflximg.com", "nflxext.com", "nflxvideo.net", "nflxso.net"}
	relaySrv     = flag.String("relay", defaultRelaySrv, "DNS server to relay requests to")
	listenAddr   = flag.String("listen", defaultListenAddr, "Address to listen for DNS requests")
	help         = flag.Bool("h", false, "Show this help message")
)

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	log.Printf("Running Netflix AAAA skipping DNS server on %s and relaying to %s", *listenAddr, *relaySrv)
	srv := startServer()

	sch := make(chan os.Signal)
	signal.Notify(sch, syscall.SIGTERM, syscall.SIGINT)
	<-sch

	srv.Shutdown()
}

func startServer() *dns.Server {
	srv := &dns.Server{Addr: *listenAddr, Net: "udp", Handler: newMuxedHandler()}
	if e := srv.ListenAndServe(); e != nil {
		log.Fatalf("Failed to start %s, error: %s", os.Args[0], e.Error())
	}
	return srv
}

func newMuxedHandler() *dns.ServeMux {
	mux := dns.NewServeMux()
	mux.HandleFunc(".", defaultHandler)
	for _, n := range netflixNames {
		mux.HandleFunc(n, netflixHandler)
	}
	return mux
}

func defaultHandler(rw dns.ResponseWriter, m *dns.Msg) {
	r, e := dns.Exchange(m, *relaySrv)
	if e != nil {
		r = new(dns.Msg)
		r.SetRcode(m, dns.RcodeServerFailure)
	}
	rw.WriteMsg(r)
}

func netflixHandler(rw dns.ResponseWriter, m *dns.Msg) {
	if m.Question[0].Qtype == dns.TypeAAAA {
		var r dns.Msg
		rw.WriteMsg(r.SetReply(m))
	} else {
		defaultHandler(rw, m)
	}
}
