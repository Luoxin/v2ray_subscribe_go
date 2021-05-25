package dns

import (
	"net"
	"strings"

	"github.com/Luoxin/Eutamias/utils"
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"
)

type DnsClient interface {
	Init() error
	LookupHost(domain string) pie.Strings
}

//go:generate pie DnsClientList.*
type DnsClientList []DnsClient

func ParseDnsUrl(dnsServiceAddr string) DnsClient {
	if utils.IsIpV4(dnsServiceAddr) {
		return NewUdpClient(net.JoinHostPort(dnsServiceAddr, "53"))
	} else if strings.HasPrefix(dnsServiceAddr, "tls://") {
		client := NewDotClient(strings.TrimPrefix(dnsServiceAddr, "tls://"))
		err := client.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil
		}
		return client
	} else if strings.HasPrefix(dnsServiceAddr, "https://") {
		client := NewDohClient(dnsServiceAddr)
		err := client.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil
		}
		return client
	} else {
		log.Warnf("unsupported %v", dnsServiceAddr)
	}

	return nil
}
