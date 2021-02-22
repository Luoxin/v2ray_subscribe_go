package proxy

import (
	"bufio"
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
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		switch {
		case strings.HasPrefix(scanner.Text(), "ssr://"):
			p, err = ParseSSRLink(strings.TrimSpace(scanner.Text()))
			if err != nil {
				log.Errorf("err:%v", err)
			}
			return

		case strings.HasPrefix(scanner.Text(), "vmess://"):
			p, err = ParseVmessLink(strings.TrimSpace(scanner.Text()))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			return

		case strings.HasPrefix(scanner.Text(), "ss://"):
			p, err = ParseSSLink(strings.TrimSpace(scanner.Text()))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			return

		case strings.HasPrefix(scanner.Text(), "trojan://"):
			p, err = ParseTrojanLink(strings.TrimSpace(scanner.Text()))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
			return

		}

	}

	err = errors.New("nonsupport content")
	return
}
