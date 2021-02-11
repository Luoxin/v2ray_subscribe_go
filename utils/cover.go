package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"

	"subsrcibe/domain"
	"subsrcibe/geolite"
	"subsrcibe/proxy"
	"subsrcibe/title"
)

type clashInfo struct {
	name string
	host string

	proxyInfo proxy.Proxy
}

type CoverSubscribe struct {
	nodeMap map[string]*clashInfo
}

func NewCoverSubscribe() *CoverSubscribe {
	return &CoverSubscribe{
		nodeMap: map[string]*clashInfo{},
	}
}

func (c *CoverSubscribe) Nodes2Clash(nodes []*domain.ProxyNode) string {

	titleGen := title.NewProxyTitle()

	for _, node := range nodes {
		if node.NodeDetail == nil {
			return ""
		}

		proxyConfig, err := proxy.ParseProxy(node.NodeDetail.Buf)
		if err != nil {
			log.Errorf("err:%v", err)
			continue
		}

		title := titleGen.Get()

		c.nodeMap[title] = &clashInfo{
			name: title,
			host: proxyConfig.BaseInfo().Server,

			proxyInfo: proxyConfig,
		}
	}

	return c.genClashConfig()
}

func (c *CoverSubscribe) genClashConfig() string {
	var nodeList, nameList []string

	var usaNameList, hkNameList, twNameList, rsNameList, jpNameList, skNameList []string

	// TODO for 生成数据
	for _, x := range c.nodeMap {
		record := geolite.GetCountry(x.host)
		if record != nil {
			if len(record.Country.Names["zh-CN"]) > 0 {
				x.name = fmt.Sprintf("(%s)%s", record.Country.Names["zh-CN"], x.name)
			}

			switch record.Country.Names["zh-CN"] {
			case "美国":
				usaNameList = append(usaNameList, x.name)
			case "香港":
				hkNameList = append(hkNameList, x.name)
			case "台湾":
				twNameList = append(twNameList, x.name)
			case "新加坡":
				rsNameList = append(rsNameList, x.name)
			case "日本":
				jpNameList = append(jpNameList, x.name)
			case "韩国":
				skNameList = append(skNameList, x.name)
			}
		}

		x.proxyInfo.SetName(x.name)

		// replace := func() string {
		// 	s := strings.ReplaceAll(x.proxyInfo.ToClash(), "\"", "")
		// 	s = strings.ReplaceAll(s, ":", ": ")
		// 	s = strings.ReplaceAll(s, ",", ", ")
		// 	return s
		// }

		nodeList = append(nodeList, x.proxyInfo.ToClash())
		nameList = append(nameList, x.name)
	}

	nodeData := map[string]interface{}{
		"ProxyNodeList": fmt.Sprintf("- %+v", strings.Join(nodeList, "\n  - ")),
		"Proxies":       fmt.Sprintf("- %+v", strings.Join(nodeList, "\n  - ")),
		"NameList":      fmt.Sprintf("- %+v", strings.Join(nameList, "\n      - ")),
	}

	if len(usaNameList) > 0 {
		nodeData["UsaNameList"] = fmt.Sprintf("- %+v", strings.Join(usaNameList, "\n      - "))
	} else {
		nodeData["UsaNameList"] = ""
	}

	if len(hkNameList) > 0 {
		nodeData["HkNameList"] = fmt.Sprintf("- %+v", strings.Join(hkNameList, "\n      - "))
	} else {
		nodeData["HkNameList"] = ""
	}

	if len(twNameList) > 0 {
		nodeData["TwNameList"] = fmt.Sprintf("- %+v", strings.Join(twNameList, "\n      - "))
	} else {
		nodeData["TwNameList"] = ""
	}

	if len(rsNameList) > 0 {
		nodeData["RsNameList"] = fmt.Sprintf("- %+v", strings.Join(rsNameList, "\n      - "))
	} else {
		nodeData["RsNameList"] = ""
	}

	if len(jpNameList) > 0 {
		nodeData["JpNameList"] = fmt.Sprintf("- %+v", strings.Join(jpNameList, "\n      - "))
	} else {
		nodeData["JpNameList"] = ""
	}

	if len(skNameList) > 0 {
		nodeData["SkNameList"] = fmt.Sprintf("- %+v", strings.Join(skNameList, "\n      - "))
	} else {
		nodeData["SkNameList"] = ""
	}

	var b bytes.Buffer
	t, err := template.New("").Parse(ClashTpl)
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}

	err = t.Execute(&b, nodeData)
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}

	return b.String()
}

func (c *CoverSubscribe) vmess2Clash(s string) domain.ClashVmess {
	s = strings.TrimPrefix(s, "vmess://")

	config, err := Base64DecodeStripped(s)
	if err != nil {
		return domain.ClashVmess{}
	}

	var vmess domain.Vmess
	err = json.Unmarshal([]byte(config), &vmess)
	if err != nil {
		log.Errorf("err:%v", err)
		return domain.ClashVmess{}
	}

	clashVmess := domain.ClashVmess{}
	clashVmess.Name = vmess.PS

	clashVmess.Type = "vmess"
	clashVmess.Server = vmess.Add
	switch vmess.Port.(type) {
	case string:
		clashVmess.Port, _ = vmess.Port.(string)
	case int:
		clashVmess.Port, _ = vmess.Port.(int)
	case float64:
		clashVmess.Port, _ = vmess.Port.(float64)
	default:

	}
	clashVmess.UUID = vmess.ID
	clashVmess.AlterID = vmess.Aid
	clashVmess.Cipher = vmess.Type
	if strings.EqualFold(vmess.TLS, "tls") {
		clashVmess.TLS = true
	} else {
		clashVmess.TLS = false
	}
	if "ws" == vmess.Net {
		clashVmess.Network = vmess.Net
		clashVmess.WSPATH = vmess.Path
	}

	return clashVmess
}
