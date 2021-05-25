package dns

import (
	"context"
	"net"

	"github.com/elliotchance/pie/pie"
	"github.com/ncruces/go-dns"
	log "github.com/sirupsen/logrus"
)

type DotClient struct {
	client         *net.Resolver
	dnsServiceAddr string
}

func NewDotClient(dnsService string) *DotClient {
	return &DotClient{
		dnsServiceAddr: dnsService,
	}
}

func (p *DotClient) Init() error {
	var err error
	p.client, err = dns.NewDoTResolver(p.dnsServiceAddr, dns.DoTCache())
	if err != nil {
		log.Debugf("err:%v", err)
		return err
	}

	return nil
}

func (p *DotClient) LookupHost(domain string) (hostList pie.Strings) {
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
