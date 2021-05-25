package dns

import (
	"context"
	"net"

	"github.com/elliotchance/pie/pie"
	"github.com/ncruces/go-dns"
	log "github.com/sirupsen/logrus"
)

type DohClient struct {
	client         *net.Resolver
	dnsServiceAddr string
}

func NewDohClient(dnsService string) *DohClient {
	return &DohClient{
		dnsServiceAddr: dnsService,
	}
}

func (p *DohClient) Init() error {
	var err error
	p.client, err = dns.NewDoHResolver(p.dnsServiceAddr, dns.DoHCache())
	if err != nil {
		log.Debugf("err:%v", err)
		return err
	}

	return nil
}

func (p *DohClient) LookupHost(domain string) (hostList pie.Strings) {
	var err error
	ctx := context.TODO()
	hostList, err = p.client.LookupHost(ctx, domain)
	if err != nil {
		log.Debugf("err:%v", err)
		return
	}

	log.Debugf("query from %v awser %+v", p.dnsServiceAddr, hostList.Join(","))

	return
}
