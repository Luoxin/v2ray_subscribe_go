package proxies

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"subsrcibe/proxy"
	"subsrcibe/proxycheck"
	"subsrcibe/title"
)

//go:generate pie ProxyList.*
type ProxyList []proxy.Proxy

func (ps ProxyList) Init() error {

	return nil
}

func (ps ProxyList) NameAddIndex() ProxyList {
	num := len(ps)
	for i := 0; i < num; i++ {
		ps[i].SetName(fmt.Sprintf("%s_%+02v", ps[i].BaseInfo().Name, i+1))
	}
	return ps
}

func (ps ProxyList) NameReIndex() ProxyList {
	num := len(ps)
	for i := 0; i < num; i++ {
		originName := ps[i].BaseInfo().Name
		country := strings.SplitN(originName, "_", 2)[0]
		ps[i].SetName(fmt.Sprintf("%s_%+02v", country, i+1))
	}
	return ps
}

func (ps ProxyList) Clone() ProxyList {
	result := make(ProxyList, 0, len(ps))
	for _, pp := range ps {
		if pp != nil {
			result = append(result, pp.Clone())
		}
	}
	return result
}

func (ps ProxyList) AppendWithUrl(contact string) ProxyList {
	p, err := proxy.ParseProxy(contact)
	if err != nil {
		return ps
	}

	if ps.FindFirstUsing(func(value proxy.Proxy) bool {
		return value.BaseInfo().GetUrl() == p.BaseInfo().GetUrl()
	}) >= 0 {
		return ps
	}

	for {
		name := title.Random()
		if ps.FindFirstUsing(func(value proxy.Proxy) bool {
			return name == value.BaseInfo().Name
		}) >= 0 {
			p.SetName(title.Random())
			break
		}
	}

	ps = append(ps, p)

	return ps
}

func (ps ProxyList) SortWithTest() (psn ProxyList) {
	check := proxycheck.NewProxyCheck()
	err := check.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	ps.Each(func(p proxy.Proxy) {
		err := check.AddWithClash(p.ToClash(), func(result proxycheck.Result) error {

			return nil
		})
		if err != nil {
			log.Errorf("err:%v", err)
		}
	})

	return
}
