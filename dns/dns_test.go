package dns

import (
	"testing"

	"github.com/Luoxin/Eutamias/log"
)

func TestDns_QueryIpv4One4AllDnsService(t *testing.T) {
	log.InitLog()
	dns := NewDns()
	err := dns.Init(10)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
	dns.AddServiceList(
		"1.2.4.8",
		"tls://dns.alidns.com",
		"https://223.5.5.5/dns-query",
	)

	t.Log(dns.QueryIpv4FastestBack("baidu.com"))
	t.Log(dns.QueryIpv4FastestBack("114.114.114.114"))
	t.Log(dns.QueryIpv4FastestBack("114.114.114.114"))
	t.Log(LockupDefault("baidu.com"))
	t.Log(LockupDefault("114.114.114.114"))
	t.Log(LookupAllHosts("github.com"))
	t.Log(dns.QueryIpv4FastestIp("baidu.com"))
	t.Log(LookupHostsFastestIp("baidu.com"))
	t.Log(LookupHostsFastestIp("google.com"))
}
