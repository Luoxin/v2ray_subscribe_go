package dns

import (
	"reflect"

	"github.com/bluele/gcache"
	"github.com/elliotchance/pie/pie"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type UdpClient struct {
	client         *dns.Client
	dnsServiceAddr string
	cache          gcache.Cache
}

// net default: udp
func NewUdpClient(dnsService string) *UdpClient {
	c := new(dns.Client)

	return &UdpClient{
		client:         c,
		dnsServiceAddr: dnsService,
		cache:          gcache.New(20).LRU().Build(),
	}
}

func (p *UdpClient) Init() error {
	return nil
}

func (p *UdpClient) LookupHost(domain string) (hostList pie.Strings) {
	cacheValue, err := p.cache.Get(domain)
	if err != nil {
		if err != gcache.KeyNotFoundError {
			log.Debugf("err:%v", err)
			return
		}
	} else if cacheValue != nil {
		hostList = cacheValue.(pie.Strings)
		if len(hostList) > 0 {
			return cacheValue.(pie.Strings)
		}
	}

	m := new(dns.Msg)

	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	r, ttl, err := p.client.Exchange(m, p.dnsServiceAddr)
	if err != nil {
		log.Debugf("err:%v", err)
		return
	}

	if r.Rcode != dns.RcodeSuccess {
		log.Debugf("invalid answer name %s after MX query for %s", domain, p.dnsServiceAddr)
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

	err = p.cache.SetWithExpire(domain, hostList, ttl)
	if err != nil {
		log.Debugf("err:%v", err)
	}

	return
}
