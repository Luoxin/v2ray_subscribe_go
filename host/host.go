package host

import (
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type HostSubRule struct {
	SubUrl       string
	LastUpdateAt time.Time
}

type Host struct {
	client *resty.Client

	hostSubRuleMap map[string]*HostSubRule
}

func NewHost() *Host {
	return &Host{
		hostSubRuleMap: map[string]*HostSubRule{},
	}
}

func (p *Host) Init() error {
	p.client = resty.New().
		SetTimeout(time.Second * 5).
		SetRetryMaxWaitTime(time.Second * 5).
		SetRetryWaitTime(time.Second).
		SetLogger(log.New())

	return nil
}
