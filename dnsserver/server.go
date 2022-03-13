package dnsserver

import (
	"context"
	"fmt"
	"log"

	"github.com/likexian/doh-go"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/yimsoijoi/stubresolver/cacher"
)

var records map[string]string

type Server struct {
	Ctx       context.Context
	DnsServer dns.Server
	DohClient doh.DoH
	Cacher     cacher.Cacher
}

type answerMap map[cacher.Key][]string

func(s *Server) parseQuery(m *dns.Msg) error {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			key := cacher.NewKey(q.Name,q.Qtype,-1)
			ip, err := s.Cacher.Get(key)
			if ip != nil && err == nil {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)		t, ok := dnsTypes[q.Qtype]
			if !ok {
				fmt.Errorf("unsupported DNS type:%v",q.Qtype)
			}
			key := cacher.NewKey(q.Name,q.Qtype,-1)
			ip, err := s.Cacher.Get(key)
			if ip != nil && err == nil {
				rr, err := dns.NewRR(q.Name,t, ip)
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		
	}
	return nil
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}
