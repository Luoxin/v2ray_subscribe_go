package proxy

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ParseProxyToClash(content string) (string, error) {
	p, err := ParseProxy(content)
	if err != nil {
		return "", err
	}

	return p.ToClash(), nil
}

func ParseProxy(content string) (p Proxy, err error) {
	content = strings.TrimSpace(content)

	switch {
	case strings.HasPrefix(content, "ssr://"):
		p, err = ParseSSRLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
		}
		return

	case strings.HasPrefix(content, "vmess://"):
		p, err = ParseVmessLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		return

	case strings.HasPrefix(content, "ss://"):
		p, err = ParseSSLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		return

	case strings.HasPrefix(content, "trojan://"):
		p, err = ParseTrojanLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
		return

	default:
		err = errors.New("nonsupport content")
	}

	return
}
