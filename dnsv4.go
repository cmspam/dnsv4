package main

import (
    "flag"
    "fmt"
    "net"
    "strings"
    "time"

    "github.com/miekg/dns"
)

var (
    upstreamServer string
    upstreamPort   int
    listenIP       string
    listenPort     int
)

// Sample DNS records (local records)
var dnsRecords = map[string]dns.RR{
    "example.org.": &dns.AAAA{
        Hdr: dns.RR_Header{Name: "example.org.", Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
        AAAA: net.ParseIP("2606:2800:220:1:248:1893:25c8:1946"),
    },
}

type CustomDNSHandler struct {
    records map[string]dns.RR
}

func NewCustomDNSHandler(records map[string]dns.RR) *CustomDNSHandler {
    return &CustomDNSHandler{
        records: records,
    }
}

func (h *CustomDNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
    m := new(dns.Msg)
    m.SetReply(r)

    question := r.Question[0]
    qName := strings.TrimSuffix(question.Name, ".")
    qType := question.Qtype

    if qType == dns.TypeA {
        // Handle A record request
        if rr, ok := h.records[qName]; ok && rr.Header().Rrtype == dns.TypeA {
            m.Answer = append(m.Answer, rr)
            w.WriteMsg(m)
            return
        }

        // If A record doesn't exist, proxy the request
        proxyReply := h.proxyQuery(r)
        if proxyReply != nil {
            w.WriteMsg(proxyReply)
        }
        return
    } else if qType == dns.TypeAAAA {
        // Handle AAAA record request
        aQuery := new(dns.Msg)
        aQuery.SetQuestion(question.Name, dns.TypeA)
        aReply := h.proxyQuery(aQuery)

        if aReply != nil {
            for _, rr := range aReply.Answer {
                if rr.Header().Rrtype == dns.TypeA {
                    // A record exists, return NOERROR with no answer
                    m.Rcode = dns.RcodeSuccess
                    w.WriteMsg(m)
                    return
                }
            }
        }

        // No A record found, check local records for AAAA
        if rr, ok := h.records[qName]; ok && rr.Header().Rrtype == dns.TypeAAAA {
            m.Answer = append(m.Answer, rr)
            w.WriteMsg(m)
            return
        }

        // If no AAAA record exists, proxy the original request
        proxyReply := h.proxyQuery(r)
        if proxyReply != nil {
            w.WriteMsg(proxyReply)
        }
        return
    }

    // Default to proxying if not an A or AAAA query
    proxyReply := h.proxyQuery(r)
    if proxyReply != nil {
        w.WriteMsg(proxyReply)
    }
}

func (h *CustomDNSHandler) proxyQuery(request *dns.Msg) *dns.Msg {
    // Proxy DNS request to an external server
    c := new(dns.Client)
    c.Timeout = time.Second
    in, _, err := c.Exchange(request, fmt.Sprintf("%s:%d", upstreamServer, upstreamPort))
    if err != nil {
        fmt.Printf("Proxy query failed: %v\n", err)
        return nil
    }
    return in
}

func main() {
    // Parse command line arguments
    flag.StringVar(&upstreamServer, "upstream-server", "1.1.1.1", "Upstream DNS server address")
    flag.IntVar(&upstreamPort, "upstream-port", 53, "Upstream DNS server port")
    flag.StringVar(&listenIP, "listen-ip", "127.0.0.1", "IP address to listen on")
    flag.IntVar(&listenPort, "listen-port", 53, "Port to listen on")
    flag.Parse()

    handler := NewCustomDNSHandler(dnsRecords)
    server := &dns.Server{
        Addr:    fmt.Sprintf("%s:%d", listenIP, listenPort),
        Net:     "udp",
        Handler: handler,
    }
    fmt.Printf("Starting DNS server on %s:%d\n", listenIP, listenPort)
    fmt.Printf("Proxying to %s:%d\n", upstreamServer, upstreamPort)
    err := server.ListenAndServe()
    if err != nil {
        fmt.Printf("Failed to start DNS server: %s\n", err)
    }
}
