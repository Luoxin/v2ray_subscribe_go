package dns

import (
	"net"
	"reflect"

	"github.com/Luoxin/Eutamias/utils"
	"github.com/elliotchance/pie/pie"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type Dns struct {
	DnsServiceList pie.Strings
}

func NewDns() *Dns {
	return &Dns{}
}

func (p *Dns) AddServiceList(serviceList ...string) *Dns {
	for _, service := range serviceList {
		if utils.IsIpV4(service) {
			service = net.JoinHostPort(service, "53")
		}

		p.DnsServiceList = append(p.DnsServiceList, service)
	}

	return p
}

func (p *Dns) QueryIpv4OneWithDnsService(dnsService, domain string) (hostList pie.Strings) {
	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, _, err := c.Exchange(m, dnsService)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Errorf("invalid answer name %s after MX query for %s", domain, dnsService)
		return
	}

	for _, a := range r.Answer {
		switch x := a.(type) {
		case *dns.A:
			hostList = append(hostList, x.A.String())
		case *dns.AAAA:
			hostList = append(hostList, x.AAAA.String())
		case *dns.CNAME:
			hostList = append(hostList, p.QueryIpv4OneWithDnsService(dnsService, x.Target)...)
		default:
			log.Errorf("unsupported type :%v", reflect.TypeOf(a))
		}
	}

	log.Debugf("query from %v awser %+v", dnsService, hostList.Join(","))

	return
}

func (p *Dns) QueryIpv4One4AllDnsService(domain string) (hostList pie.Strings) {
	p.DnsServiceList.Each(func(dnsService string) {
		hostList = append(hostList, p.QueryIpv4OneWithDnsService(dnsService, domain)...)
	})
	return
}
