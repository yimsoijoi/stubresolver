package dohclient

import (
	"context"
	"log"

	"github.com/likexian/doh-go"
	dohdns "github.com/likexian/doh-go/dns"
)

type dohCli struct {
	Dcli *doh.DoH
}

func New() *dohCli {
	client := doh.Use(doh.GoogleProvider, doh.CloudflareProvider)
	return &dohCli{
		Dcli: client,
	}
}

func (d *dohCli) GetAnswer(ctx context.Context, dName dohdns.Domain) []dohdns.Answer {
	resp, err := d.Dcli.Query(ctx, dohdns.Domain(dName), dohdns.TypeA)
	if err != nil {
		log.Println("failed to query", err.Error())
	}
	return resp.Answer
}
