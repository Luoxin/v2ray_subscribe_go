package geolite

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dreamacro/clash/component/trie"
	"github.com/Dreamacro/clash/dns"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/elliotchance/pie/pie"
	"github.com/oschwald/geoip2-golang"
	"github.com/oschwald/maxminddb-golang"
	log "github.com/sirupsen/logrus"

	country2 "github.com/Luoxin/Eutamias/country"
)

var db *geoip2.Reader
var dnsResolver *dns.Resolver

const (
	geoLiteUrl    = "http://api.luoxin.live/api/eutamias/file/GeoLite2.mmdb"
	geoLiteDbName = "GeoLite2.mmdb"
)

func InitGeoLite() error {
	geoLite2Path := filepath.Join(utils.GetConfigDir(), geoLiteDbName)

	execPath := utils.GetExecPath()
	pwdPath, _ := os.Getwd()

	var retryCount = 0
RETRY:
	retryCount++
	if utils.FileExists(filepath.Join(execPath, geoLiteDbName)) {
		err := utils.CopyFile(filepath.Join(execPath, geoLiteDbName), geoLite2Path)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	} else if utils.FileExists(filepath.Join(pwdPath, geoLiteDbName)) {
		err := utils.CopyFile(filepath.Join(pwdPath, geoLiteDbName), geoLite2Path)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	} else if !utils.FileExists(geoLite2Path) {
		log.Info("downloading geo Lite...")
		err := utils.DownloadWithProgressbar(geoLiteUrl, geoLite2Path)
		if err != nil {
			log.Errorf("err:%v", err)
			e := os.Remove(geoLite2Path)
			if e != nil {
				log.Errorf("err:%v", e)
				log.Errorf("remove %v fail. please delete manually", geoLite2Path)
				return e
			}
			return err
		}
	}

	log.Infof("loading from %v", geoLite2Path)
	var err error
	db, err = geoip2.Open(geoLite2Path)
	if err != nil {
		switch err.(type) {
		case maxminddb.InvalidDatabaseError:
			if retryCount > 4 {
				return err
			}
			e := os.Remove(geoLite2Path)
			if e != nil {
				log.Errorf("err:%v", e)
				log.Errorf("remove %v fail. please delete manually", geoLite2Path)
				return e
			}
			log.Warnf("geo lite file damage, try fix")
			goto RETRY
		}
		log.Errorf("err:%v", err)
		return err
	}

	var dnsServices = pie.Strings{
		"tls://dns.alidns.com:853",
		"tls://dns.cfiec.net:853",
		"https://2400:3200:baba::1/dns-query",
		"https://2400:3200::1/dns-query",
		"https://dns.cfiec.net/dns-query",
		"https://223.5.5.5/dns-query",
		"https://223.6.6.6/dns-query",
		"https://dns.alidns.com/dns-query",
		"https://dns.ipv6dns.com/dns-query",
		"tls://dns.ipv6dns.com:853",
		"tls://dns.pub:853",
		"tls://doh.pub:853",
		"https://doh.pub/dns-query",
		"223.5.5.5",
		"2400:3200::1",
		"2001:dc7:1000::1",
		"2400:da00::6666",
		"2001:cc0:2fff:1::6666",
		"114.114.114.114",
		"1.2.4.8",
		"180.76.76.76",
		"119.29.29.29",
		"119.28.28.28",
		"https://dns.google/dns-query",
		"tls://dns.google:853",
		"https://dns.quad9.net/dns-query",
		"https://dns11.quad9.net/dns-query",
		"https://dns.twnic.tw/dns-query",
		"https://1.1.1.1/dns-query",
		"https://1.0.0.1/dns-query",
		"https://cloudflare-dns.com/dns-query",
		"https://dns.adguard.com/dns-query",
		"https://doh.dns.sb/dns-query",
		"tls://185.184.222.222@853",
		"tls://185.222.222.222@853",
		"https://doh-jp.blahdns.com/dns-query",
		"https://public.dns.iij.jp/dns-query",
		"https://v6.rubyfish.cn/dns-query",
		"tls://v6.rubyfish.cn:853",
		"https://[2001:4860:4860::6464]/dns-query",
		"https://[2001:4860:4860::64]/dns-query",
		"https://[2606:4700:4700::1111]/dns-query",
		"https://[2606:4700:4700::64]/dns-query",
		"https://dns.quad9.net/dns-query",
		"tls://2a09::@853",
		"tls://2a09::1@853",
		"8.8.8.8",
		"203.112.2.4",
		"9.9.9.9",
		"101.101.101.101",
		"203.80.96.10",
		"218.102.23.228",
		"61.10.0.130",
		"202.181.240.44",
		"112.121.178.187",
		"168.95.192.1",
		"202.76.4.1",
		"202.14.67.4",
	}

	var nameServices []dns.NameServer
	dnsServices.Unique().Each(func(service string) {
		if strings.HasPrefix(service, "tls://") {
			nameServices = append(nameServices, dns.NameServer{
				Net:  "tcp-tls",
				Addr: strings.TrimPrefix(service, "tls://"),
			})
		} else if strings.HasPrefix(service, "http://") {
			// return
			// nameServices = append(nameServices, dns.NameServer{
			// 	Net:  "http",
			// 	Addr: service,
			// })
		} else if strings.HasPrefix(service, "https://") {
			// return
			// nameServices = append(nameServices, dns.NameServer{
			// 	Net:  "https",
			// 	Addr: service,
			// })
		} else {
			nameServices = append(nameServices, dns.NameServer{
				Net:  "",
				Addr: service + ":53",
			})
		}
	})

	dnsConfig := dns.Config{
		Main:           nameServices,
		Fallback:       nameServices,
		Default:        nameServices,
		IPv6:           true,
		EnhancedMode:   dns.MAPPING,
		FallbackFilter: dns.FallbackFilter{},
		Pool:           nil,
		Hosts:          trie.New(),
	}

	dnsResolver = dns.NewResolver(dnsConfig)

	// overtrurConfig := config.Config{
	// 	BindAddress:              "7891",
	// 	DebugHTTPAddress:         "127.0.0.1:2006",
	// 	PrimaryDNS:               []*common.DNSUpstream{},
	// 	OnlyPrimaryDNS:           true,
	// 	IPv6UseAlternativeDNS:    false,
	// 	AlternativeDNSConcurrent: true,
	// 	IPNetworkFile: struct {
	// 		Primary     string
	// 		Alternative string
	// 	}{
	// 		Primary:     "./ip_network_primary_sample",
	// 		Alternative: "./ip_network_alternative_sample",
	// 	},
	// 	DomainFile: struct {
	// 		Primary            string
	// 		Alternative        string
	// 		PrimaryMatcher     string
	// 		AlternativeMatcher string
	// 		Matcher            string
	// 	}{
	// 		Primary:            "./domain_primary",
	// 		Alternative:        "./domain_alternative",
	// 		PrimaryMatcher:     "",
	// 		AlternativeMatcher: "",
	// 		Matcher:            "full-map",
	// 	},
	// 	HostsFile: struct {
	// 		HostsFile string
	// 		Finder    string
	// 	}{
	// 		HostsFile: "./hosts",
	// 		Finder:    "full-map",
	// 	},
	// 	MinimumTTL:    86400,
	// 	DomainTTLFile: "./domain_ttl",
	// 	CacheSize:     10000,
	// 	// RejectQType:                 255,
	// 	WhenPrimaryDNSAnswerNoneUse: "primaryDNS",
	// }
	//
	// dnsServices.Unique().Each(func(service string) {
	// 	var dnsUpstream common.DNSUpstream
	// 	if strings.HasPrefix(service, "tls://") {
	// 		return
	// 		nameServices = append(nameServices, dns.NameServer{
	// 			Net:  "tcp-tls",
	// 			Addr: strings.TrimPrefix(service, "tls://"),
	// 		})
	// 	} else if strings.HasPrefix(service, "http://") {
	// 		return
	// 		nameServices = append(nameServices, dns.NameServer{
	// 			Net:  "http",
	// 			Addr: service,
	// 		})
	// 	} else if strings.HasPrefix(service, "https://") {
	// 		return
	// 		nameServices = append(nameServices, dns.NameServer{
	// 			Net:  "https",
	// 			Addr: service,
	// 		})
	// 	} else {
	// 		dnsUpstream = common.DNSUpstream{
	// 			Name:     service,
	// 			Address:  service,
	// 			Protocol: "udp",
	// 			Timeout:  1,
	// 			EDNSClientSubnet: &common.EDNSClientSubnetType{
	// 				Policy:     "disable",
	// 				ExternalIP: "",
	// 				NoCookie:   true,
	// 			},
	// 		}
	// 	}
	//
	// 	overtrurConfig.PrimaryDNS = append(overtrurConfig.PrimaryDNS, &dnsUpstream)
	// })
	//
	// dispatcher := outbound.Dispatcher{
	// 	PrimaryDNS:                  overtrurConfig.PrimaryDNS,
	// 	AlternativeDNS:              overtrurConfig.AlternativeDNS,
	// 	OnlyPrimaryDNS:              overtrurConfig.OnlyPrimaryDNS,
	// 	WhenPrimaryDNSAnswerNoneUse: overtrurConfig.WhenPrimaryDNSAnswerNoneUse,
	// 	// IPNetworkPrimarySet:         overtrurConfig.IPNetworkPrimarySet,
	// 	// IPNetworkAlternativeSet:     overtrurConfig.IPNetworkAlternativeSet,
	// 	DomainPrimaryList:     overtrurConfig.DomainPrimaryList,
	// 	DomainAlternativeList: overtrurConfig.DomainAlternativeList,
	//
	// 	RedirectIPv6Record:       overtrurConfig.IPv6UseAlternativeDNS,
	// 	AlternativeDNSConcurrent: overtrurConfig.AlternativeDNSConcurrent,
	// 	MinimumTTL:               overtrurConfig.MinimumTTL,
	// 	DomainTTLMap:             overtrurConfig.DomainTTLMap,
	//
	// 	Hosts: overtrurConfig.Hosts,
	// 	Cache: overtrurConfig.Cache,
	// }
	// dispatcher.Init()

	// ns, err := dnsResolver.ResolveIP("google.com")
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	//
	// log.Info(ns)

	return nil
}

type Country struct {
	Code   string
	CnName string
	EnName string
	Emoji  string
}

func GetCountry(host string) (c *Country, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic:%v", r)
		}
	}()

	ns, err := dnsResolver.ResolveIPv4(host)
	if err != nil {
		return nil, err
	}
	ip := ns.String()

	// ns, err := net.LookupIP(host)
	// if err != nil {
	// 	log.Errorf("err:%v", err)
	// 	return nil, err
	// }
	//
	// if len(ns) == 0 {
	// 	return nil, errors.New("not found ip")
	// }
	//
	// ip := ns[0].String()

	record, err := db.City(net.ParseIP(ip))
	if err != nil {
		return nil, err
	}

	c = &Country{
		CnName: record.Country.Names["zh-CN"],
		EnName: record.Country.Names["en"],
	}

	country := country2.Country.GetByEnName(strings.ReplaceAll(record.Country.Names["en"], " ", ""))
	if country == nil {
		country = country2.Country.GetByCnName(record.Country.Names["zh-CN"])
	}

	if country != nil {
		c.Code = country.Code
		c.Emoji = country.Emoji
		c.CnName = country.CnName
		c.EnName = country.EnName
	} else {
		log.Warnf("not found country(host:%s, ip:%s):%+v", host, ip, record)
	}

	return c, nil
}
