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
	DnsServer *dns.Server
	DohClient *dohclient.DohCli
	Cacher    *cacher.Cacher
}

func New(ctx context.Context, dnsserver *dns.Server, dohcli *dohclient.DohCli, cacher *cacher.Cacher) *Server {
	return &Server{
		Ctx:       ctx,
		DnsServer: dnsserver,
		DohClient: dohcli,
		Cacher:    cacher,
	}
}

func NewDNSServer() *dns.Server {
	return &dns.Server{
		Addr: "127.0.0.1:5300",
		Net:  "udp",
	}
}

func NewRR(dName, t, v string) (dns.RR, error) {
	return dns.NewRR(fmt.Sprintf("%s %s %s", dName, t, v))
}

type typeRRMap map[uint16]string

var rrMap = typeRRMap{
	1:  "A",
	28: "AAAA",
	5:  "CNAME",
}

func (s *Server) Worker(m *dns.Msg) error {
	for _, q := range m.Question {
		t, ok := rrMap[q.Qtype]
		if !ok {
			continue
		}
		k := cacher.NewKey(q.Name, t, -1)
		v, err := s.Cacher.Get(k)
		if err != nil {
			log.Println("Cache missed")
		}
		if v != "" {
			// found cache
			rr, err := NewRR(q.Name, t, v)
			if err != nil {
				return errors.Wrapf(err, "NewRR failed key: %s, type: %s, value: %s", q.Name, t, v)
			}
			m.Answer = append(m.Answer, rr)
		} else {
			// DoH
			answers, err := s.DohClient.GetAnswer(s.Ctx, dohdns.Domain(q.Name))
			if err != nil {
				return err
			}
			for _, answer := range answers {
				t, ok := rrMap[uint16(answer.Type)]
				if !ok {
					log.Println("unsportted type", t)
					continue
				}
				rr, err := NewRR(answer.Name, t, answer.Data)
				if err != nil {
					return errors.Wrapf(err, "NewRR failed key: %s, type: %s, value: %s", answer.Name, t, answer.Data)
				}
				key := cacher.NewKey(answer.Name, t, answer.TTL)
				s.Cacher.Set(key, answer.Data, answer.TTL)
				m.Answer = append(m.Answer, rr)
			}
		}
	}
	return nil
}

func (s *Server) HandleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		s.Worker(m)
	}

	w.WriteMsg(m)
}
