# NETFLIX DNS - Skip netflix AAAA DNS requests

netflix-dns is a DNS server that skips AAAA requests for netflix domains. This is helpful for people who use tunnelbroker.net IPV6 tunnels.

## How to use

## Simple

Start netflix-dns,
```
# netflix-dns
2017/11/16 16:27:54 Running Netflix AAAA skipping DNS server on :53 and relaying to 8.8.8.8:53
```

This starts a DNS server listening on the default DNS port on your system.

Set DNS server for you system to `127.0.0.1`, and you're done.

netflix-dns relays DNS requests to the Google public DNS server (8.8.8.8). A different server can be set using the `-relay` flag.
```
# netflix-dns -relay 208.67.222.123:53
2017/11/16 16:28:35 Running Netflix AAAA skipping DNS server on :53 and relaying to 208.67.222.123:53
```

### With NetworkManager and dnsmasq

Systems which already have a local dnsmasq resolver with NetworkManager need to start netflix-dns on a different listening port.
A different port can be set using the `-listen` flag.

```
# netflix-dns -listen :2053
2017/11/16 16:29:04 Running Netflix AAAA skipping DNS server on :2053 and relaying to 8.8.8.8:53
```

Configure NetworkManager and dnsmasq to forward DNS requests for netflix domains to netflix-dns.
```
# cat /etc/NetworkManager/dnsmasq.d/10-netflix-dns
server=/netflix.com/127.0.0.1#2053
server=/netflix.net/127.0.0.1#2053
server=/nflxext.com/127.0.0.1#2053
server=/nflximg.com/127.0.0.1#2053
server=/nflxvideo.net/127.0.0.1#2053
server=/nflxso.net/127.0.0.1#2053
```
Restart NetworkManager after creating the file.

