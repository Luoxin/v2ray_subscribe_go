package proxies

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/geolite"
	"github.com/Luoxin/Eutamias/proxy"
	"github.com/Luoxin/Eutamias/proxycheck"
	"github.com/Luoxin/Eutamias/title"
)

//go:generate pie ProxyList.*
type ProxyList []proxy.Proxy

type Proxies struct {
	proxyList  ProxyList
	proxyTitle *title.ProxyTitle

	proxyMap  map[string]bool
	proxyLock sync.Mutex

	netEaseProxyList ProxyList
}

func NewProxies() *Proxies {
	return &Proxies{
		proxyList:  ProxyList{},
		proxyMap:   make(map[string]bool),
		proxyTitle: title.NewProxyTitle(),
	}
}

func (ps *Proxies) Len() int {
	return len(ps.proxyList)
}

func (ps *Proxies) NameAddIndex() *Proxies {
	num := len(ps.proxyList)
	for i := 0; i < num; i++ {
		ps.proxyList[i].SetName(fmt.Sprintf("%s_%+02v", ps.proxyList[i].BaseInfo().Name, i+1))
	}
	return ps
}

func (ps *Proxies) NameReIndex() *Proxies {
	num := len(ps.proxyList)
	for i := 0; i < num; i++ {
		originName := ps.proxyList[i].BaseInfo().Name
		country := strings.SplitN(originName, "_", 2)[0]
		ps.proxyList[i].SetName(fmt.Sprintf("%s_%+02v", country, i+1))
	}
	return ps
}

func (ps *Proxies) Clone() *Proxies {
	result := make(ProxyList, 0, len(ps.proxyList))
	for _, pp := range ps.proxyList {
		if pp != nil {
			result = append(result, pp.Clone())
		}
	}
	return &Proxies{proxyList: result}
}

func (ps *Proxies) AppendWithUrl(contact string) *Proxies {
	p, err := proxy.ParseProxy(contact)
	if err != nil {
		return ps
	}

	p.SetName("Proxies")

	baseUrl := p.Link()

	// 去重
	{
		ps.proxyLock.Lock()
		exist := ps.proxyMap[baseUrl]
		ps.proxyLock.Unlock()

		if exist {
			return nil
		}
	}

	// 改名字
	p.SetName(ps.proxyTitle.Get())

	c, err := geolite.GetCountry(p.BaseInfo().Server)
	if err == nil {
		p.SetCountry(c.CnName)
		p.SetName(fmt.Sprintf("(%s)%s", c.CnName, p.BaseInfo().Name))
		p.SetEmoji(c.Emoji)
	}

	ps.proxyLock.Lock()
	ps.proxyMap[baseUrl] = true
	ps.proxyLock.Unlock()

	ps.proxyList = append(ps.proxyList, p)
	return ps
}

func (ps *Proxies) Append(contact string, name string) *Proxies {
	p, err := proxy.ParseProxy(contact)
	if err != nil {
		return ps
	}

	p.SetName("Proxies")

	baseUrl := p.Link()

	// 去重
	{
		ps.proxyLock.Lock()
		exist := ps.proxyMap[baseUrl]
		ps.proxyLock.Unlock()

		if exist {
			return nil
		}
	}

	// 改名字
	p.SetName(name)

	c, err := geolite.GetCountry(p.BaseInfo().Server)
	if err == nil {
		p.SetCountry(c.CnName)
		p.SetName(fmt.Sprintf("(%s)%s", c.CnName, p.BaseInfo().Name))
		p.SetEmoji(c.Emoji)
	}

	ps.proxyLock.Lock()
	ps.proxyMap[baseUrl] = true
	ps.proxyLock.Unlock()

	ps.proxyList = append(ps.proxyList, p)
	return ps
}

func (ps *Proxies) AppendNetEaseWithUrl(contact string) *Proxies {
	p, err := proxy.ParseProxy(contact)
	if err != nil {
		return ps
	}

	p.SetName(fmt.Sprintf("网易云%x", md5.Sum([]byte(contact))))

	ps.netEaseProxyList = append(ps.netEaseProxyList, p)
	return ps
}

func (ps *Proxies) GetUsableList() (psn *Proxies) {
	check := proxycheck.NewProxyCheck()
	err := check.Init(0)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	psn = NewProxies()
	ps.proxyList.Each(func(p proxy.Proxy) {
		err := check.AddWithLink(p.Link(), func(result proxycheck.Result) (err error) {
			if result.Err != nil {
				return nil
			}

			if result.Delay < 0 || result.Speed < 0 {
				return nil
			}

			if result.Speed > 5000 {
				return nil
			}

			psn.Append(result.ProxyUrl, func() string {
				if p.BaseInfo() == nil {
					return ps.proxyTitle.Get()
				}

				return p.BaseInfo().Name
			}())
			return nil
		})
		if err != nil {
			log.Errorf("err:%v", err)
		}
	})

	check.WaitFinish()

	if psn.Len() == 0 {
		return ps
	}

	return psn
}

type countryNode struct {
	Name, Emoji string
	NameList    pie.Strings
	TestUrl     string
}

type DnsType uint32

const (
	DnsTypeNil DnsType = iota
	DnsTypeBase
	DnsTypeLocal
)

type ClashConfig struct {
	DnsType DnsType
}

func (ps *Proxies) ToClashConfig(c ClashConfig) string {
	var proxyList, proxyNameList, countryGroupList []string

	countryMap := map[string]*countryNode{
		// "香港": {
		// 	Name:     "香港",
		// 	Emoji:    "🇭🇰",
		// 	NameList: []string{},
		// },
		"台湾省": {
			Name:     "台湾省",
			Emoji:    "🇹🇼",
			NameList: []string{},
		},
	}

	ps.proxyList.Each(func(p proxy.Proxy) {
		proxyList = append(proxyList, p.ToClash())
		proxyNameList = append(proxyNameList, p.BaseInfo().Name)

		if p.BaseInfo().Country != "" {
			if _, ok := countryMap[p.BaseInfo().Country]; !ok {
				countryMap[p.BaseInfo().Country] = &countryNode{
					Name:     p.BaseInfo().Country,
					Emoji:    p.BaseInfo().Emoji,
					NameList: []string{},
				}
				countryGroupList = append(countryGroupList, fmt.Sprintf("%s %s", p.BaseInfo().Emoji, p.BaseInfo().Country))
			}

			countryMap[p.BaseInfo().Country].NameList = append(countryMap[p.BaseInfo().Country].NameList, p.BaseInfo().Name)
		}

	})

	var countryNodeList []*countryNode
	for _, v := range countryMap {
		v.NameList = v.NameList.Unique()

		if len(v.NameList) == 0 {
			v.NameList = append(v.NameList, "DIRECT")
		}

		v.TestUrl = "http://www.gstatic.com/generate_204"
		countryNodeList = append(countryNodeList, v)
	}

	var netEaseProxyList, netEaseProxyNameList []string
	{
		for _, p := range ps.netEaseProxyList {
			netEaseProxyList = append(netEaseProxyList, p.ToClash())
			netEaseProxyNameList = append(netEaseProxyNameList, p.BaseInfo().Name)
		}
	}

	val := &clashTplVal{
		Dns:                  baseDns,
		ProxyList:            proxyList,
		ProxyNameList:        proxyNameList,
		CountryNodeList:      countryNodeList,
		CountryGroupList:     countryGroupList,
		NetEaseProxyList:     netEaseProxyList,
		NetEaseProxyNameList: netEaseProxyNameList,
		TestUrl:              "http://www.gstatic.com/generate_204",
		MixedPort:            conf.Config.Proxy.MixedPort,
	}

	switch c.DnsType {
	case DnsTypeLocal:
		if conf.Config.Dns.EnableService {
			t, err := template.New("").Parse(dnsTpl)
			if err != nil {
				log.Errorf("err:%v", err)
				return ""
			}

			var b bytes.Buffer
			err = t.Execute(&b, map[string]interface{}{
				"DnsServiceList": []string{
					fmt.Sprintf("127.0.0.1:%d", conf.Config.Dns.ServicePort),
				},
			})
			if err != nil {
				log.Errorf("err:%v", err)
				return ""
			} else {
				val.Dns = b.String()
			}
		} else {
			log.Warnf("local dns service not enable")
		}
	}

	_lock.RLock()
	t, err := template.New("").Parse(clashTpl)
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}
	_lock.RUnlock()

	var b bytes.Buffer
	err = t.Execute(&b, val)
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}

	return b.String()
}

func (ps *Proxies) ToV2rayConfig() string {
	var linkList pie.Strings
	ps.proxyList.Each(func(p proxy.Proxy) {
		linkList = append(linkList, p.Link())
	})

	return base64.RawStdEncoding.EncodeToString([]byte(strings.Join(linkList, "\n")))
}
