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

	case strings.HasPrefix(content, "vmess://"):
		p, err = ParseVmessLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	case strings.HasPrefix(content, "ss://"):
		p, err = ParseSSLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	case strings.HasPrefix(content, "trojan://"):
		p, err = ParseTrojanLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	case strings.HasPrefix(content, "http://"):
		p, err = ParseHttpLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	case strings.HasPrefix(content, "socket://"):
		p, err = ParseSocketLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	case strings.HasPrefix(content, "socket4://"):
		p, err = ParseSocketLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	case strings.HasPrefix(content, "socket5://"):
		p, err = ParseSocketLink(content)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	default:
		err = errors.New("nonsupport content")
	}

	if strings.HasPrefix(p.BaseInfo().Server, "127.0.0.1") {
		err = errors.New("host is local")
	}

	return
}
