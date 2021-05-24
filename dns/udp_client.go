package dns

import (
	"reflect"

	"github.com/elliotchance/pie/pie"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type UdpClient struct {
	client         *dns.Client
	dnsServiceAddr string
}

func NewUdpClient(dnsService string) *UdpClient {
	c := new(dns.Client)

	return &UdpClient{
		client:         c,
		dnsServiceAddr: dnsService,
	}
}

func (p *UdpClient) Init() error {
	return nil
}

func (p *UdpClient) LookupHost(domain string) (hostList pie.Strings) {

	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := p.client.Exchange(m, p.dnsServiceAddr)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Errorf("invalid answer name %s after MX query for %s", domain, p.dnsServiceAddr)
		return
	}

	for _, a := range r.Answer {
		switch x := a.(type) {
		case *dns.A:
			hostList = append(hostList, x.A.String())
		case *dns.AAAA:
			hostList = append(hostList, x.AAAA.String())
		case *dns.CNAME:
			hostList = append(hostList, p.LookupHost(x.Target)...)
		default:
			log.Errorf("unsupported type :%v", reflect.TypeOf(a))
		}
	}

	log.Debugf("query from %v awser %+v", p.dnsServiceAddr, hostList.Join(","))

	return
}
