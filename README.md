# dnsv4

A DNS proxy server that prioritizes IPv4 over IPv6 connections while maintaining IPv6 fallback support. This helps prevent connectivity issues in environments where IPv6 may be unreliable or slower than IPv4, or in any other situation where IPv4 is preferred.

## Purpose

dnsv4 is designed to optimize connectivity by prioritizing IPv4 over IPv6 at the DNS level while still allowing IPv6 fallback. It does this by:

1. When a client requests an AAAA (IPv6) record, dnsv4 first checks if an A (IPv4) record exists
2. If an IPv4 record exists, the AAAA query returns empty (suppressing IPv6)
3. If no IPv4 record exists, the AAAA query proceeds normally, allowing IPv6 connectivity

This approach ensures:
- Applications will try IPv4 first when both are available
- IPv6 still works for domains that are IPv6-only
- Better reliability in environments with problematic IPv6 connectivity
- IPv6 timeout delays are avoided

## Prerequisites

- Go 1.18 or higher
- `github.com/miekg/dns` package

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/cmspam/dnsv4
    cd dnsv4
    ```

2. Install dependencies:
    ```bash
    go mod init dnsv4
    go get github.com/miekg/dns
    ```

3. Build the binary:
    ```bash
    go build
    ```

## Usage

Run the server with default settings:

```bash
sudo ./dnsv4
```

The default configuration:
- Listens on `127.0.0.1:53`
- Uses Cloudflare's DNS (`1.1.1.1:53`) as upstream server

### Command Line Options

- `-listen`: Address and port to listen on (default: "127.0.0.1:53")
- `-upstream`: Upstream DNS server address and port (default: "1.1.1.1:53")

Example with custom settings:

```bash
sudo ./dnsv4 -listen "0.0.0.0:5353" -upstream "8.8.8.8:53"
```

Note: Running on port 53 requires root/administrator privileges

