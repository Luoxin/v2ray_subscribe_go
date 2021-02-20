package proxies

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"

	"subscribe/geolite"
	"subscribe/proxy"
	"subscribe/proxycheck"
	"subscribe/title"
)

//go:generate pie ProxyList.*
type ProxyList []proxy.Proxy

type Proxies struct {
	proxyList ProxyList
	nameList  pie.Strings
}

func NewProxies() *Proxies {
	return &Proxies{
		proxyList: ProxyList{},
	}
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

	// 去重
	if ps.proxyList.FindFirstUsing(func(value proxy.Proxy) bool {
		return value.BaseInfo().GetUrl() == p.BaseInfo().GetUrl()
	}) >= 0 {
		return ps
	}

	// 改名字
	for {
		name := title.Random()

		if !ps.nameList.Contains(name) {
			ps.nameList = append(ps.nameList, name)
			p.SetName(name)
			break
		}
	}

	c, err := geolite.GetCountry(p.BaseInfo().Server)
	if err == nil {
		p.SetCountry(c.CnName)
		p.SetName(fmt.Sprintf("(%s)%s", c.CnName, p.BaseInfo().Name))
		p.SetEmoji(c.Emoji)
	}

	ps.proxyList = append(ps.proxyList, p)
	return ps
}

func (ps *Proxies) SortWithTest() (psn *Proxies) {
	check := proxycheck.NewProxyCheck()
	err := check.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	ps.proxyList.Each(func(p proxy.Proxy) {
		err := check.AddWithClash(p.ToClash(), func(result proxycheck.Result) error {

			return nil
		})
		if err != nil {
			log.Errorf("err:%v", err)
		}
	})

	return
}

func (ps *Proxies) ToClashConfig() string {
	var proxyList, proxyNameList, countryGroupList []string

	type countryNode struct {
		Name, Emoji string
		NameList    pie.Strings
		TestUrl     string
	}

	countryMap := map[string]*countryNode{
		"香港": {
			Name:     "香港",
			Emoji:    "🇭🇰",
			NameList: []string{},
		},
		"台湾": {
			Name:     "台湾",
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

	var b bytes.Buffer
	t, err := template.New("").Parse(ClashTpl)
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}

	err = t.Execute(&b, map[string]interface{}{
		"ProxyList":        proxyList,
		"ProxyNameList":    proxyNameList,
		"CountryNodeList":  countryNodeList,
		"CountryGroupList": countryGroupList,
		"TestUrl":          "http://www.gstatic.com/generate_204",
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}

	return b.String()
}
