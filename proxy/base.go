package proxy

import (
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

/* Base implements interface Proxy. It's the basic proxy struct. Vmess etc extends Base*/
type Base struct {
	Name    string `yaml:"name" json:"name"`
	Server  string `yaml:"server" json:"server"`
	Path    string `yaml:"path" json:"path,omitempty"`
	Type    string `yaml:"type" json:"type"`
	Country string `yaml:"country,omitempty" json:"country,omitempty"`
	Emoji   string `yaml:"-" json:"-"`
	Port    int    `yaml:"port" json:"port" `
	UDP     bool   `yaml:"udp,omitempty" json:"udp,omitempty"`
	Useable bool   `yaml:"useable,omitempty" json:"useable,omitempty"`
}

func (b *Base) TypeName() string {
	if b.Type == "" {
		return "unknown"
	}
	return b.Type
}

func (b *Base) SetName(name string) {
	b.Name = name
}

func (b *Base) AddToName(name string) {
	b.Name = b.Name + name
}

func (b *Base) SetIP(ip string) {
	b.Server = ip
}

func (b *Base) BaseInfo() *Base {
	return b
}

func (b *Base) GetUrl() string {
	return fmt.Sprintf("%s:%d:%d", b.Server, b.Port, b.Port)
}

func (b *Base) Clone() Base {
	c := *b
	return c
}

func (b *Base) SetUseable(useable bool) {
	b.Useable = useable
}

func (b *Base) SetCountry(country string) {
	b.Country = country
}

func (b *Base) SetEmoji(emoji string) {
	b.Emoji = emoji
}

type Proxy interface {
	String() string
	ToClash() string
	ToSurge() string
	Link() string
	Identifier() string
	SetName(name string)
	AddToName(name string)
	SetIP(ip string)
	TypeName() string //ss ssr vmess trojan
	BaseInfo() *Base
	Clone() Proxy
	SetUseable(useable bool)
	SetCountry(country string)
	SetEmoji(emoji string)
}

func ParseProxyFromClashProxy(p map[string]interface{}) (proxy Proxy, err error) {
	pjson, err := json.Marshal(p)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	switch p["type"].(string) {
	case "ss":
		var proxy Shadowsocks
		err := json.Unmarshal(pjson, &proxy)
		if err != nil {
			return nil, err
		}
		return &proxy, nil
	case "ssr":
		var proxy ShadowsocksR
		err := json.Unmarshal(pjson, &proxy)
		if err != nil {
			return nil, err
		}
		return &proxy, nil
	case "vmess":
		var proxy Vmess
		err := json.Unmarshal(pjson, &proxy)
		if err != nil {
			return nil, err
		}
		return &proxy, nil
	case "trojan":
		var proxy Trojan
		err := json.Unmarshal(pjson, &proxy)
		if err != nil {
			return nil, err
		}
		return &proxy, nil
	case "http":
		var proxy Http
		err := json.Unmarshal(pjson, &proxy)
		if err != nil {
			return nil, err
		}
		return &proxy, nil
	case "socket":
		var proxy Socket
		err := json.Unmarshal(pjson, &proxy)
		if err != nil {
			return nil, err
		}
		return &proxy, nil
	}

	return nil, errors.New("clash json parse failed")
}
