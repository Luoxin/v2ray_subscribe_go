package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Luoxin/Eutamias/proxy"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type FuzzyMatchingParser struct {
}

func NewFuzzyMatchingParser() *FuzzyMatchingParser {
	return &FuzzyMatchingParser{}
}

func (n FuzzyMatchingParser) ParserText(body string) pie.Strings {
	var nodeList pie.Strings

	adds := func(nodeUrls ...string) {
		nodeList = append(nodeList, nodeUrls...)
	}

	add := func(nodeUrl string) {
		adds(nodeUrl)
	}

	reg := func(str string) {
		nodeLL := regexp.MustCompile(str).
			FindAllStringSubmatch(body, -1)
		for _, ll := range nodeLL {
			pie.
				Strings(ll).
				Each(add)

		}
	}

	reg(`vmess://[^\s]*`)
	reg(`trojan://[^\s]*`)
	reg(`ssr://[^\s]*`)
	reg(`ss://[^\s]*`)
	reg(`vless://[^\s]*`)

	base64List := regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$`).
		FindAllStringSubmatch(body, -1)

	for _, x := range base64List {
		pie.Strings(x).Each(func(s string) {
			str, err := utils.Base64DecodeStripped(s)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			adds(strings.Split(str, "\n")...)
		})
	}

	type clashConfig struct {
		Proxy []map[interface{}]interface{} `yaml:"proxies"`
	}

	var clashConf clashConfig
	err := yaml.Unmarshal([]byte(body), &clashConf)
	if err == nil {
		var convert func(m map[interface{}]interface{}) map[string]interface{}
		convert = func(m map[interface{}]interface{}) map[string]interface{} {
			res := map[string]interface{}{}
			for k, v := range m {
				switch v2 := v.(type) {
				case map[interface{}]interface{}:
					res[fmt.Sprint(k)] = convert(v2)
				default:
					res[fmt.Sprint(k)] = v
				}
			}
			return res
		}

		for _, x := range clashConf.Proxy {
			p, err := proxy.ParseProxyFromClashProxy(convert(x))
			if err == nil {
				add(p.Link())
			}

		}

	}

	return nodeList
}
