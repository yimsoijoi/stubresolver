package dohclient

import (
	"context"

	"github.com/likexian/doh-go"
	dohdns "github.com/likexian/doh-go/dns"
	"github.com/pkg/errors"
)

type DohCli struct {
	Dcli *doh.DoH
}

func New() *DohCli {
	client := doh.Use(doh.GoogleProvider, doh.CloudflareProvider)
	return &DohCli{
		Dcli: client,
	}
}

func (d *DohCli) GetAnswer(ctx context.Context, dName dohdns.Domain) ([]dohdns.Answer, error) {
	resp, err := d.Dcli.Query(ctx, dohdns.Domain(dName), dohdns.TypeANY)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get response from DoHdns", dName)
	}
	return resp.Answer, nil
}
