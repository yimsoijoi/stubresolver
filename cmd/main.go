package main

import (
	"context"
	"log"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"github.com/yimsoijoi/stubresolver/cacher"
	"github.com/yimsoijoi/stubresolver/config"
	"github.com/yimsoijoi/stubresolver/dnsserver"
	"github.com/yimsoijoi/stubresolver/dohclient"
	"github.com/yimsoijoi/stubresolver/rediswrapper"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		errors.Wrap(err, "failed to load config")
	}
	ctx := context.Background()
	redisCli := rediswrapper.New(ctx, conf.RedisConfig)
	cacher := cacher.New(redisCli)
	dnsServer := dnsserver.NewDNSServer(conf.DnsserverConfig)
	dohClient := dohclient.New()
	server := dnsserver.New(ctx, dnsServer, dohClient, cacher)

	dns.HandleFunc(".", server.HandleDnsRequest)
	log.Println("Starting stubborn DNS resolver on localhost:5300")
	if err := server.DnsServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// client := doh.Use(doh.GoogleProvider, doh.CloudflareProvider)

// resp, err := client.Query(ctx, dohdns.Domain(dom), dohdns.TypeA)
// if err != nil {
// 	log.Fatalln("failed to query", err.Error())
// }
// answer := resp.Answer

// rdb := *redis.NewClient(&redis.Options{
// 	DB: 1,
// })
// answerJson, err := json.Marshal(answer)
// if err != nil {
// 	log.Fatalln("failed to marshal answer", err.Error())
// }

// shortest := answer[0]
// for _, a := range resp.Answer {
// 	if a.TTL > shortest.TTL {
// 		shortest = a
// 	}
// }
// expire := time.Duration(shortest.TTL)
// if _, err := rdb.Set(ctx, dom, string(answerJson), expire).Result(); err != nil {
// 	log.Fatalln("failed to set redis", dom, answerJson)
// }

// val, err := rdb.Get(ctx, dom).Result()
// if err != nil {
// 	log.Fatalln("failed to get from redis", err.Error())
// }

// dns.HandleFunc("service.", dns.NewRR)
