# dnsv4
DNS Proxy to Prioritize IPV4

This was quickly thrown together using AI chatbots.

The idea is a DNS proxy which will force IPV4, but allow IPV6 if IPV4 is not available.

It will do the following:

If an A record exists, return the A record, and do not return the AAAA record.

If an A record doesn't exist, proxy the request unmodified.

This way, we can ensure we only get ipv4 records if they are available.
