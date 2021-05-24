package dns

import (
	"github.com/elliotchance/pie/pie"
)

type DnsClient interface {
	Init() error
	LookupHost(domain string) pie.Strings
}
