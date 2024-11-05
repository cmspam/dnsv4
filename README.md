# dnsv4
A DNS proxy server that prioritizes IPv4 over IPv6 connections while maintaining IPv6 fallback support. This helps prevent connectivity issues in environments where IPv6 may be unreliable or slower than IPv4.

## Purpose
dnsv4 is designed to optimize connectivity by prioritizing IPv4 over IPv6 while still allowing IPv6 fallback. It does this by:

- When a client requests an AAAA (IPv6) record, dnsv4 first checks if an A (IPv4) record exists
- If an IPv4 record exists, the AAAA query returns empty (suppressing IPv6)
- If no IPv4 record exists, the AAAA query proceeds normally, allowing IPv6 connectivity

This approach ensures:

- Applications will try IPv4 first when both are available
- IPv6 still works for domains that are IPv6-only
- Better reliability in environments with problematic IPv6 connectivity
- Faster connections by avoiding IPv6 timeout delays when IPv4 is available

## Prerequisites

Go 1.18 or higher
github.com/miekg/dns package

Installation

Clone the repository:
```
git clone https://github.com/cmspam/dnsv4
cd dnsv4
Install dependencies:
go mod init dnsv4
go get github.com/miekg/dns
Build the binary:
go build
```

## Configuration
dnsv4 can be configured through command-line flags or a JSON configuration file.
Command Line Options
```
-listen: Address and port to listen on (default: "127.0.0.1:53")
-upstream: Upstream DNS server address and port (default: "1.1.1.1:53")
-proxy-only: Disable IPv4 prioritization (default: false)
-config: Path to configuration file
```

## Configuration File
Create a JSON file (e.g., config.json):
```
{
    "upstream": "1.1.1.1:53",
    "listen": "127.0.0.1:53",
    "proxy_only": false
}
```
## Usage

Basic usage with defaults:
``` sudo ./dnsv4 ```

Using command line options:

``` sudo ./dnsv4 -upstream "8.8.8.8:53" -listen "0.0.0.0:5353" ```

Using a config file:

``` sudo ./dnsv4 -config config.json ```

Enable proxy-only mode (disable IPv4 prioritization):

``` sudo ./dnsv4 -proxy-only ```

Note: Running on port 53 requires root/administrator privileges
The default configuration:

- Listens on 127.0.0.1:53
- Uses Cloudflare's DNS (1.1.1.1:53) as upstream server
- IPv4 prioritization enabled

## How It Works

For regular DNS queries (A records, MX records, etc.), dnsv4 acts as a transparent proxy
For AAAA (IPv6) queries (when not in proxy-only mode):

- First checks if an A record exists for the domain
- If an A record exists, returns an empty success response
- If no A record exists, proxies the AAAA query to get IPv6 records


In proxy-only mode, all queries are forwarded to the upstream server without modification

Testing
You can test the DNS server using dig or nslookup:

Test IPv4 prioritization (when not in proxy-only mode)

``` dig @127.0.0.1 example.com AAAA    # Should return empty if A record exists ```

Test IPv6 fallback

``` dig @127.0.0.1 ipv6.google.com AAAA    # Should return IPv6 if no A record ```

Test regular IPv4 resolution

``` dig @127.0.0.1 example.com A ```

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.
