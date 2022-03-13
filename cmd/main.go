package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/likexian/doh-go"
	dohdns "github.com/likexian/doh-go/dns"
)

func main() {
	if len(os.Args) < 2 {
		panic("not enough argument")
	}
	dom := os.Args[1]
	ctx := context.Background()

	client := doh.Use(doh.GoogleProvider, doh.CloudflareProvider)

	resp, err := client.Query(ctx, dohdns.Domain(dom), dohdns.TypeA)
	if err != nil {
		log.Fatalln("failed to query", err.Error())
	}
	answer := resp.Answer

	rdb := *redis.NewClient(&redis.Options{
		DB: 1,
	})
	answerJson, err := json.Marshal(answer)
	if err != nil {
		log.Fatalln("failed to marshal answer", err.Error())
	}

	shortest := answer[0]
	for _, a := range resp.Answer {
		if a.TTL > shortest.TTL {
			shortest = a
		}
	}
	expire := time.Duration(shortest.TTL)
	if _, err := rdb.Set(ctx, dom, string(answerJson), expire).Result(); err != nil {
		log.Fatalln("failed to set redis", dom, answerJson)
	}

	// val, err := rdb.Get(ctx, dom).Result()
	// if err != nil {
	// 	log.Fatalln("failed to get from redis", err.Error())
	// }

	// dns.HandleFunc("service.", dns.NewRR)

}
