package proxy

import (
	"errors"
)

var ErrorTypeCanNotConvert = errors.New("type not support")

// Convert2SS convert proxy to ShadowsocksR if possible
func Convert2SSR(p Proxy) (ssr *ShadowsocksR, err error) {
	if p.TypeName() == "ss" {
		ss := p.(*Shadowsocks)
		if ss == nil {
			return nil, errors.New("ss is nil")
		}
		base := ss.Base
		base.Type = "ssr"
		return &ShadowsocksR{
			Base:     base,
			Password: ss.Password,
			Cipher:   ss.Cipher,
			Protocol: "origin",
			Obfs:     "plain",
			Group:    "",
		}, nil
	}
	return nil, ErrorTypeCanNotConvert
}

// Convert2SS convert proxy to Shadowsocks if possible
func Convert2SS(p Proxy) (ss *Shadowsocks, err error) {
	if p.TypeName() == "ss" {
		ssr := p.(*ShadowsocksR)
		if ssr == nil {
			return nil, errors.New("ssr is nil")
		}
		if ssr.Protocol != "origin" || ssr.Obfs != "plain" {
			return nil, errors.New("protocol or obfs not allowed")
		}
		base := ssr.Base
		base.Type = "ss"
		return &Shadowsocks{
			Base:       base,
			Password:   ssr.Password,
			Cipher:     ssr.Cipher,
			Plugin:     "",
			PluginOpts: nil,
		}, nil
	}
	return nil, ErrorTypeCanNotConvert
}
