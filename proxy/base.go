package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
)

/* Base implements interface Proxy. It's the basic proxy struct. Vmess etc extends Base*/
type Base struct {
	Name    string `yaml:"name" json:"name" gorm:"index"`
	Server  string `yaml:"server" json:"server" gorm:"index"`
	Path    string `yaml:"path" json:"path,omitempty" gorm:"index"`
	Type    string `yaml:"type" json:"type" gorm:"index"`
	Country string `yaml:"country,omitempty" json:"country,omitempty" gorm:"index"`
	Port    int    `yaml:"port" json:"port" gorm:"index"`
	UDP     bool   `yaml:"udp,omitempty" json:"udp,omitempty"`
	Useable bool   `yaml:"useable,omitempty" json:"useable,omitempty" gorm:"index"`
}

// TypeName() Get specific proxy type
func (b *Base) TypeName() string {
	if b.Type == "" {
		return "unknown"
	}
	return b.Type
}

// SetName() to a proxy
func (b *Base) SetName(name string) {
	b.Name = name
}

func (b *Base) AddToName(name string) {
	b.Name = b.Name + name
}

// SetIP() to a proxy
func (b *Base) SetIP(ip string) {
	b.Server = ip
}

// BaseInfo() get basic info struct of a proxy
func (b *Base) BaseInfo() *Base {
	return b
}

func (b *Base) GetUrl() string {
	return fmt.Sprintf("%s:%d:%s", b.Server, b.Port, b.Port)
}

// Clone() returns a new basic proxy
func (b *Base) Clone() Base {
	c := *b
	return c
}

// SetUseable() set Base info "Useable" (true or false)
func (b *Base) SetUseable(useable bool) {
	b.Useable = useable
}

// SetUseable() set Base info "Country" (string)
func (b *Base) SetCountry(country string) {
	b.Country = country
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
}

func ParseProxyFromClashProxy(p map[string]interface{}) (proxy Proxy, err error) {
	p["name"] = ""
	pjson, err := json.Marshal(p)
	if err != nil {
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
	}
	return nil, errors.New("clash json parse failed")
}
