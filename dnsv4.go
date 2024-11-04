package main

import (
    "flag"
    "fmt"
    "net"
    "time"

    "github.com/miekg/dns"
)

var (
    upstreamServer = flag.String("upstream", "1.1.1.1:53", "Upstream DNS server")
    listen         = flag.String("listen", "127.0.0.1:53", "Listen address")
)

type handler struct {
    records map[string]dns.RR
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
    m := new(dns.Msg)
    m.SetReply(r)

    q := r.Question[0]
    if rr, ok := h.records[q.Name]; ok && rr.Header().Rrtype == q.Qtype {
        m.Answer = append(m.Answer, rr)
        w.WriteMsg(m)
        return
    }

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

    records := map[string]dns.RR{
        "example.org.": &dns.AAAA{
            Hdr:  dns.RR_Header{Name: "example.org.", Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
            AAAA: net.ParseIP("2606:2800:220:1:248:1893:25c8:1946"),
        },
    }

    server := &dns.Server{
        Addr:    *listen,
        Net:     "udp",
        Handler: &handler{records: records},
    }

    fmt.Printf("Starting DNS server on %s, proxying to %s\n", *listen, *upstreamServer)
    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("Failed to start server: %s\n", err)
    }
}
