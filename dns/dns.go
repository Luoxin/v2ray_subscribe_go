package dns

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/alitto/pond"
	"github.com/bluele/gcache"
	"github.com/elliotchance/pie/pie"
	"github.com/go-ping/ping"
	"github.com/miekg/dns"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

var fastCache = gcache.New(1024).LRU().Build()
var pool = pond.New(100, 100)

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

// QueryIpv4FastestBack 最快返回的结果
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

// QueryIpv4FastestIp 最快的ip
func (p *Dns) QueryIpv4FastestIp(domain string) string {
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
			_ = fastCache.SetWithExpire(domain, fastestIp, time.Minute*30)
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

// func CheckServer(){
// 	timeout := time.Duration(5 * time.Second)
// 	t1 := time.Now()
// 	_, err := net.DialTimeout("tcp","www.google.com:443", timeout)
// 	fmt.Println("waist time :", time.Now().Sub(t1))
// 	if err != nil {
// 		fmt.Println("Site unreachable, error: ", err)
// 		return
// 	}
// 	fmt.Println("tcp server is ok")
// }

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
	if utils.IsIp(domain) {
		return domain
	}

	var host string
	val, err := fastCache.Get(domain)
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
	if utils.IsIp(domain) {
		return domain
	}

	if dnsClient == nil {
		return FastestIp(domain, LockupDefault(domain))
	} else {
		return dnsClient.QueryIpv4FastestIp(domain)
	}
}

func InitDnsService() error {
	if !conf.Config.Dns.EnableService {
		log.Warnf("dns service not start")
		return nil
	}

	server := &dns.Server{Addr: ":" + strconv.Itoa(int(conf.Config.Dns.ServicePort)), Net: "udp"}
	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:
			for _, q := range m.Question {
				switch q.Qtype {
				case dns.TypeA:
					log.Debugf("lookup %v", q.Name)
					// ipList := LookupAllHosts(q.Name)
					// ipList.Each(func(ip string) {
					// 	rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
					// 	if err == nil {
					// 		m.Answer = append(m.Answer, rr)
					// 	}
					// 	log.Infof("[dns query]%v %v", q.Name, ip)
					// })

					ip := LookupHostsFastestBack(q.Name)
					if ip != "" {
						rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
						if err == nil {
							m.Answer = append(m.Answer, rr)
						}
						log.Infof("[dns query]%v %v", q.Name, ip)
					}
				}
			}
		default:
			log.Warnf("unsupported opcode")
		}

		err := w.WriteMsg(m)
		if err != nil {
			log.Errorf("err:%v", err)
		}
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}()
	go func() {
		<-sigCh
		log.Info("dns service stop")
		err := server.Shutdown()
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}()

	return nil
}
