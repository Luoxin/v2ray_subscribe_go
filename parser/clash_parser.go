package parser

import (
	"github.com/Dreamacro/clash/config"
	"github.com/Luoxin/Eutamias/proxy"
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"
)

type ClashParser struct {
}

func (c ClashParser) ParserText(body string) pie.Strings {
	clashConf, err := config.UnmarshalRawConfig([]byte(body))
	if err != nil {
		log.Errorf("err:%v", err)
		return nil
	}

	var links pie.Strings
	for _, proxyInfo := range clashConf.Proxy {
		p, err := proxy.ParseProxyFromClashProxy(proxyInfo)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}

		links = append(links, p.Link())
	}

	return links
}

func NewClashParser() *ClashParser {
	return &ClashParser{}
}
