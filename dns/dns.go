package dns

import (
	"net"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/alitto/pond"
	"github.com/bluele/gcache"
	"github.com/elliotchance/pie/pie"
	"github.com/go-ping/ping"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

var dnsCache = gcache.New(1024).LRU().Build()
var pool = pond.New(100, 100)

// Dns TODO ipv6的支持
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

func (p *Dns) QueryIpV4(domain string) (hostList pie.Strings) {
	if domain == "" {
		return pie.Strings{}
	}

	p.dnsClientList.Each(func(client DnsClient) {
		hostList = append(hostList, client.LookupHost(domain)...)
	})
	hostList = hostList.Unique()
	return
}

// QueryIpv4FastestBack 最快返回的结果
func (p *Dns) QueryIpv4FastestBack(domain string) string {
	if domain == "" {
		return ""
	}

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

// QueryIpv4FastestIp 最快的ip
func (p *Dns) QueryIpv4FastestIp(domain string) string {
	if domain == "" {
		return ""
	}

	var hostList pie.Strings
	p.dnsClientList.Each(func(client DnsClient) {
		hostList = append(hostList, client.LookupHost(domain)...)
	})

	return FastestIp(domain, hostList)
}

func FastestIp(domain string, ipList pie.Strings) string {
	var delay = time.Duration(-1)
	fastestIp := ""
	ipList.Unique().Each(func(ip string) {
		d := Ping(ip)
		if d < 0 {
			return
		}

		if d < delay || delay < 0 {
			delay = d
			fastestIp = ip
			_ = dnsCache.SetWithExpire(domain, fastestIp, time.Minute)
		}
	})

	if fastestIp == "" {
		fastestIp = ipList.First()
	}

	return fastestIp
}

func Ping(ip string) time.Duration {
	i := ping.New(ip)
	i.Count = 5
	i.Interval = time.Millisecond * 10
	i.Timeout = time.Millisecond * 4500
	i.Debug = conf.Config.Debug
	i.Size = 128

	i.SetLogger(log.New())

	i.SetPrivileged(true)

	i.OnRecv = func(packet *ping.Packet) {
		log.Debugf("%v,rtt:%v", packet.Addr, packet.Rtt)
	}

	i.OnDuplicateRecv = func(packet *ping.Packet) {
		log.Debugf("%v,rtt:%v", packet.Addr, packet.Rtt)
	}

	i.OnFinish = func(statistics *ping.Statistics) {
		log.Debugf("%v,maxRtt:%v,minRtt:%v,avgRtt:%v", statistics.Addr, statistics.MaxRtt, statistics.MinRtt, statistics.AvgRtt)
	}

	err := i.Run()
	if err != nil {
		log.Debugf("err:%v", err)
		return -1
	}

	return i.Statistics().AvgRtt
}

var dnsClient *Dns

type InitClient struct {
}

func (p *InitClient) Init() (needRun bool, err error) {
	return InitDnsClient()
}

func (p *InitClient) WaitFinish() {
}

func (p *InitClient) Name() string {
	return "dns client"
}

func InitDnsClient() (bool, error) {
	if !conf.Config.Dns.Enable {
		log.Debugf("dns client not start")
		return false, nil
	}

	dnsClient = NewDns()
	err := dnsClient.Init(0)
	if err != nil {
		log.Errorf("err:%v", err)
		return false, err
	}

	dnsClient.AddServiceList(conf.Config.Dns.Nameserver...)

	return true, nil
}

func LockupDefault(domain string) (hostList pie.Strings) {
	if domain == "" {
		return pie.Strings{}
	}

	if utils.IsIp(domain) {
		hostList = append(hostList, domain)
		return
	}

	hostList, err := net.LookupHost(domain)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	return hostList
}

func LookupHostsFastestBack(domain string) string {
	if domain == "" {
		return ""
	}

	if utils.IsIp(domain) {
		return domain
	}

	var host string
	val, err := dnsCache.Get(domain)
	if err == nil {
		host = val.(string)
	} else {
		if dnsClient == nil {
			host = LockupDefault(domain).First()
		} else {
			host = dnsClient.QueryIpv4FastestBack(domain)
		}

		pool.Submit(func() {
			LookupHostsFastestIp(domain)
		})
	}

	return host
}

func LookupAllHosts(domain string) (hostList pie.Strings) {
	if domain == "" {
		return pie.Strings{}
	}

	if utils.IsIp(domain) {
		hostList = append(hostList, domain)
		return
	}

	if dnsClient == nil {
		return LockupDefault(domain)
	}

	return dnsClient.QueryIpV4(domain)
}

func LookupHostsFastestIp(domain string) string {
	if domain == "" {
		return ""
	}

	if utils.IsIp(domain) {
		return domain
	}

	if dnsClient == nil {
		return FastestIp(domain, LockupDefault(domain))
	} else {
		return dnsClient.QueryIpv4FastestIp(domain)
	}
}
