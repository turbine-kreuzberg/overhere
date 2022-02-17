# Overhere

This is a very minimal DNS server. \
It only supports IPv4 and IPv6 entries.

It uses the default upstream to resolve requests, \
if the resolution fails is responds with its own IP.

## Help

``` bash
overhere -h
NAME:
   Overhere - A very minimal DNS server for development purposes.

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --resolve-to value, --ip value  IP to resolve to. Autodetected by default. [$OVERHERE_RESOLVE_TO]
   --port value, -p value          Port to listen on. (default: 15353) [$OVERHERE_PORT]
   --verbose, -v                   Verbosity of logging. (default: false) [$OVERHERE_VERBOSE]
   --help, -h                      show help (default: false)
```

## Example

Run server:

``` bash
go run .
```

Resolve development domain:
``` bash
dig @127.0.0.1 -p 15353 my-development-server

; <<>> DiG 9.16.1-Ubuntu <<>> @127.0.0.1 -p 15353 my-development-server
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 14706
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;my-development-server.         IN      A

;; ANSWER SECTION:
my-development-server.  60      IN      A       192.168.1.2

;; Query time: 11 msec
;; SERVER: 127.0.0.1#15353(127.0.0.1)
;; WHEN: Do Feb 17 09:08:52 CET 2022
;; MSG SIZE  rcvd: 76
```

Resolve google.com domain:
``` bash
dig @127.0.0.1 -p 15353 google.com           

; <<>> DiG 9.16.1-Ubuntu <<>> @127.0.0.1 -p 15353 google.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 47675
;; flags: qr aa rd; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;google.com.                    IN      A

;; ANSWER SECTION:
google.com.             60      IN      A       142.251.36.206
google.com.             60      IN      AAAA    2a00:1450:4016:809::200e

;; Query time: 15 msec
;; SERVER: 127.0.0.1#15353(127.0.0.1)
;; WHEN: Do Feb 17 09:10:01 CET 2022
;; MSG SIZE  rcvd: 92
```
