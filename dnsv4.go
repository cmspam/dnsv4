package main

import (
    "flag"
    "fmt"
    "time"
    "encoding/json"
    "os"

    "github.com/miekg/dns"
)

type Config struct {
    UpstreamServer string `json:"upstream"`
    ListenAddress  string `json:"listen"`
    ProxyOnly      bool   `json:"proxy_only"`
}

var (
    upstreamServer = flag.String("upstream", "1.1.1.1:53", "Upstream DNS server")
    listen         = flag.String("listen", "127.0.0.1:53", "Listen address")
    proxyOnly      = flag.Bool("proxy-only", false, "Disable IPv4 prioritization")
    configFile     = flag.String("config", "", "Path to config file")
)

type handler struct {
    config *Config
}

func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
    m := new(dns.Msg)
    m.SetReply(r)

    q := r.Question[0]
    
    // For AAAA queries, check if A record exists (unless in proxy-only mode)
    if !h.config.ProxyOnly && q.Qtype == dns.TypeAAAA {
        aQuery := new(dns.Msg)
        aQuery.SetQuestion(q.Name, dns.TypeA)
        if reply, _, err := new(dns.Client).Exchange(aQuery, h.config.UpstreamServer); err == nil && len(reply.Answer) > 0 {
            w.WriteMsg(m) // Return empty success if A record exists
            return
        }
    }

    // Proxy the query
    if reply, _, err := (&dns.Client{Timeout: time.Second}).Exchange(r, h.config.UpstreamServer); err == nil {
        w.WriteMsg(reply)
    }
}

func loadConfig(path string) (*Config, error) {
    // Default config
    config := &Config{
        UpstreamServer: "1.1.1.1:53",
        ListenAddress:  "127.0.0.1:53",
        ProxyOnly:      false,
    }

    // If config file specified, load it
    if path != "" {
        file, err := os.ReadFile(path)
        if err != nil {
            return nil, fmt.Errorf("error reading config file: %v", err)
        }

        if err := json.Unmarshal(file, config); err != nil {
            return nil, fmt.Errorf("error parsing config file: %v", err)
        }
    }

    // Command line flags override config file
    flag.Visit(func(f *flag.Flag) {
        switch f.Name {
        case "upstream":
            config.UpstreamServer = *upstreamServer
        case "listen":
            config.ListenAddress = *listen
        case "proxy-only":
            config.ProxyOnly = *proxyOnly
        }
    })

    return config, nil
}

func main() {
    flag.Parse()

    config, err := loadConfig(*configFile)
    if err != nil {
        fmt.Printf("Failed to load config: %s\n", err)
        os.Exit(1)
    }

    server := &dns.Server{
        Addr:    config.ListenAddress,
        Net:     "udp",
        Handler: &handler{config: config},
    }

    fmt.Printf("Starting DNS server on %s, proxying to %s\n", config.ListenAddress, config.UpstreamServer)
    if config.ProxyOnly {
        fmt.Println("Running in proxy-only mode (IPv4 prioritization disabled)")
    }

    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("Failed to start server: %s\n", err)
    }
}
