package parser

import (
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"subsrcibe/utils"
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

	pie.
		Strings(regexp.MustCompile(`vmess://[^\s]*`).
			FindStringSubmatch(body)).
		Each(add)

	pie.
		Strings(regexp.MustCompile(`trojan://[^\s]*`).
			FindStringSubmatch(body)).
		Each(add)

	pie.
		Strings(regexp.MustCompile(`ssr://[^\s]*`).
			FindStringSubmatch(body)).
		Each(add)

	pie.
		Strings(regexp.MustCompile(`ss://[^\s]*`).
			FindStringSubmatch(body)).
		Each(add)

	pie.
		Strings(regexp.MustCompile(`vless://[^\s]*`).
			FindStringSubmatch(body)).
		Each(add)

	pie.
		Strings(regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$`).
			FindStringSubmatch(body)).
		Each(func(s string) {
			str, err := utils.Base64DecodeStripped(s)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			adds(strings.Split(str, "\n")...)
		})

	return nodeList
}
