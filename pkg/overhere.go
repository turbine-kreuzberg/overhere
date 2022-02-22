package overhere

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/miekg/dns"
)

// NewServer sets up a new DNS server.
func NewServer(defaultIP string, addr string, port int, verbose bool) (*dns.Server, error) {
	ip := net.ParseIP(defaultIP)
	var err error

	if ip == nil {
		ip, err = GetOutboundIP()
		if err != nil {
			log.Fatalf("lookup outbound IP: %v", err)
		}

		log.Printf("autodetected IP: %s", ip)
	}

	log.Printf("resolving to IP: %s", ip)

	srv := &dns.Server{
		Addr: addr + ":" + strconv.Itoa(port),
		Net:  "udp",
	}
	srv.Handler = &handler{
		defaultIP: ip,
		verbose:   verbose,
	}

	return srv, nil
}

type handler struct {
	defaultIP net.IP
	verbose   bool
}

// ServeDNS resolves dns requests.
// It forwards the request upstream and proxies existing domains.
// If the upstream request fails (domains is unknown) it answers with the defaultIP.
func (h *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		switch question.Qtype {
		case dns.TypeA:
			domain := question.Name

			if h.verbose {
				log.Printf("resolving %s", domain)
			}

			ips, err := net.LookupIP(domain)
			if err == nil {
				if h.verbose {
					log.Printf("resolved %s to %s", domain, ips)
				}

				for _, ip := range ips {
					addAnswer(&msg, domain, ip)
				}
				continue
			}

			//ip := net.ParseIP(w.LocalAddr().String())
			ip := h.defaultIP

			if h.verbose {
				log.Printf("returning default IP (%s)", ip)
			}

			addAnswer(&msg, domain, ip)
		}
	}

	err := w.WriteMsg(&msg)
	if err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func addAnswer(msg *dns.Msg, domain string, ip net.IP) {
	if ip.To4() == nil {
		answer := &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   domain,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			AAAA: ip,
		}
		msg.Answer = append(msg.Answer, answer)
	} else {
		answer := &dns.A{
			Hdr: dns.RR_Header{
				Name:   domain,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			A: ip,
		}
		msg.Answer = append(msg.Answer, answer)
	}
}

// GetOutboundIP opens a udp connection to detect the "primary" IP
func GetOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	ip := net.ParseIP(localAddr[0:idx])
	return ip, nil
}
