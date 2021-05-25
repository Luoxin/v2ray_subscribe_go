package dns

import (
	"net"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/elliotchance/pie/pie"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

// TODO ipv6的支持
type Dns struct {
	dnsClientList DnsClientList
	pool          *ants.Pool
}

func NewDns() *Dns {
	return &Dns{}
}

func (p *Dns) Init(size int) error {
	if size == 0 {
		size = 10
	}

	var err error
	p.pool, err = ants.NewPool(size)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (p *Dns) AddServiceList(serviceList ...string) *Dns {
	for _, service := range serviceList {
		client := ParseDnsUrl(service)
		if client != nil {
			log.Debugf("add dns %v", service)
			p.dnsClientList = append(p.dnsClientList, client)
		} else {
			log.Warnf("add dns service %v fail", service)
		}
	}

	return p
}

// func (p *Dns) QueryIpv4OneWithDnsService(dnsService, domain string) (hostList pie.Strings) {
// 	c := new(dns.Client)
// 	m := new(dns.Msg)
//
// 	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
// 	m.RecursionDesired = true
//
// 	r, _, err := c.Exchange(m, dnsService)
// 	if err != nil {
// 		log.Errorf("err:%v", err)
// 		return
// 	}
//
// 	if r.Rcode != dns.RcodeSuccess {
// 		log.Errorf("invalid answer name %s after MX query for %s", domain, dnsService)
// 		return
// 	}
//
// 	for _, a := range r.Answer {
// 		switch x := a.(type) {
// 		case *dns.A:
// 			hostList = append(hostList, x.A.String())
// 		case *dns.AAAA:
// 			hostList = append(hostList, x.AAAA.String())
// 		case *dns.CNAME:
// 			hostList = append(hostList, p.QueryIpv4OneWithDnsService(dnsService, x.Target)...)
// 		default:
// 			log.Errorf("unsupported type :%v", reflect.TypeOf(a))
// 		}
// 	}
//
// 	log.Debugf("query from %v awser %+v", dnsService, hostList.Join(","))
//
// 	return
// }

func (p *Dns) QueryIpV4(domain string) (hostList pie.Strings) {
	p.dnsClientList.Each(func(client DnsClient) {
		hostList = append(hostList, client.LookupHost(domain)...)
	})
	hostList = hostList.Unique()
	return
}

func (p *Dns) QueryIpv4FastestBack(domain string) string {
	c := make(chan string)

	go func() {
		p.dnsClientList.Each(func(client DnsClient) {
			go func() {
				hostList := client.LookupHost(domain)
				host := hostList.FirstOr("")
				if host == "" {
					return
				}

				select {
				case c <- host:
				case <-time.After(time.Millisecond * 10):
				}
			}()
		})
	}()

	return <-c
}

var dnsClient *Dns

func InitDnsClient() error {
	if !conf.Config.Dns.Enable {
		log.Warnf("dns client not start")
		return nil
	}

	dnsClient = NewDns()
	err := dnsClient.Init(0)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	dnsClient.AddServiceList(conf.Config.Dns.Nameserver...)

	return nil
}

func LockupDefault(domain string) (hostList pie.Strings) {
	hostList, err := net.LookupHost(domain)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	return hostList
}

func LookupHostsFastestBack(domain string) string {
	if dnsClient == nil {
		return LockupDefault(domain).First()
	}

	return dnsClient.QueryIpv4FastestBack(domain)
}

func LookupAllHosts(domain string) (hostList pie.Strings) {
	if dnsClient == nil {
		return LockupDefault(domain)
	}

	return dnsClient.QueryIpV4(domain)
}
