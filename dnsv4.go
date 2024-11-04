package main

import (
    "flag"
    "fmt"
    "time"

    "github.com/miekg/dns"
)

var (
    upstreamServer = flag.String("upstream", "1.1.1.1:53", "Upstream DNS server")
    listen         = flag.String("listen", "127.0.0.1:53", "Listen address")
)

type handler struct{}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
    m := new(dns.Msg)
    m.SetReply(r)

    q := r.Question[0]
    
    // For AAAA queries, check if A record exists
    if q.Qtype == dns.TypeAAAA {
        aQuery := new(dns.Msg)
        aQuery.SetQuestion(q.Name, dns.TypeA)
        if reply, _, err := new(dns.Client).Exchange(aQuery, *upstreamServer); err == nil && len(reply.Answer) > 0 {
            w.WriteMsg(m) // Return empty success if A record exists
            return
        }
    }

    // Proxy the query
    if reply, _, err := (&dns.Client{Timeout: time.Second}).Exchange(r, *upstreamServer); err == nil {
        w.WriteMsg(reply)
    }
}

func main() {
    flag.Parse()

    server := &dns.Server{
        Addr:    *listen,
        Net:     "udp",
        Handler: &handler{},
    }

    fmt.Printf("Starting DNS server on %s, proxying to %s\n", *listen, *upstreamServer)
    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("Failed to start server: %s\n", err)
    }
}
