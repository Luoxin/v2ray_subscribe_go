package dns

import (
	"testing"

	"github.com/Luoxin/Eutamias/log"
)

func TestDns_QueryIpv4One4AllDnsService(t *testing.T) {
	log.InitLog()
	dns := NewDns()
	dns.AddServiceList("114.114.114.114:53",
		"223.5.5.5",
		"8.8.8.8",
		"1.2.4.8",
		"119.28.28.28",
		"180.76.76.76",
		"tls://dns.alidns.com:853")
	dns.QueryIpv4One4AllDnsService("baidu.com")
}
