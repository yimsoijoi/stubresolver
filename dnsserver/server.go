package dnsserver

import (
	"context"
	"fmt"
	"log"

	dohdns "github.com/likexian/doh-go/dns"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/yimsoijoi/stubresolver/cacher"
	"github.com/yimsoijoi/stubresolver/dohclient"
)

type Server struct {
	Ctx       context.Context
	DnsServer dns.Server
	DohClient dohclient.DohCli
	Cacher    cacher.Cacher
}

type answerMap map[cacher.Key][]string

// func (s *Server) RR(dName, t, v string) (dns.RR, error) {
// 	return dns.NewRR(fmt.Sprintf("%s %s %d", dName, t, v))
// }

func (s *Server) Worker(m *dns.Msg) error {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("query for %s (TypeA)\n ", q.Name)
			k := cacher.NewKey(q.Name, "A", -1)
			ip, err := s.Cacher.Get(k)
			if err == nil {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
				return err
			}
		case dns.TypeAAAA:
			log.Printf("query for %s (TypeAAAA)\n ", q.Name)
			k := cacher.NewKey(q.Name, "AAAA", -1)
			ip, err := s.Cacher.Get(k)
			if err == nil {
				rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
				return err
			}
		default:
			log.Println("DoH query")
			dom := dohdns.Domain(q.Name)
			dohAnswers, err := s.DohClient.GetAnswer(s.Ctx, dom)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to get DoH answer for %s", dom))
			}
			for _, ans := range dohAnswers {
				switch ans.Type {
				case int(dns.TypeA):
					log.Printf("query for %s (TypeA)\n ", q.Name)
					k := cacher.NewKey(q.Name, "A", -1)
					ip, err := s.Cacher.Get(k)
					if err == nil {
						rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
						return err
					}
				case int(dns.TypeAAAA):
					log.Printf("query for %s (TypeAAAA)\n ", q.Name)
					k := cacher.NewKey(q.Name, "AAAA", -1)
					ip, err := s.Cacher.Get(k)
					if err == nil {
						rr, err := dns.NewRR(fmt.Sprintf("%s AAAA %s", q.Name, ip))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
						return err
					}
				}

			}

		}
	}
	return nil
}

// // func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
// 	m := new(dns.Msg)
// 	m.SetReply(r)
// 	m.Compress = false

// 	switch r.Opcode {
// 	case dns.OpcodeQuery:
// 		parseQuery(m)
// 	}

// 	w.WriteMsg(m)
// }
