package utils

import (
	"strings"

	"github.com/Luoxin/Eutamias/domain"
)

func GetProxyNodeType(u string) domain.ProxyNodeType {
	if strings.HasPrefix(u, "vmess://") {
		return domain.ProxyNodeTypeVmess
	} else if strings.HasPrefix(u, "trojan://") {
		return domain.ProxyNodeTypeTrojan
	} else if strings.HasPrefix(u, "vless://") {
		return domain.ProxyNodeTypeVless
	} else if strings.HasPrefix(u, "ssr://") {
		return domain.ProxyNodeTypeSSR
	} else if strings.HasPrefix(u, "ss://") {
		return domain.ProxyNodeTypeSS
	} else if strings.HasPrefix(u, "http://") {
		return domain.ProxyNodeTypeHttp
	} else if strings.HasPrefix(u, "socket://") {
		return domain.ProxyNodeTypeSocket
	} else if strings.HasPrefix(u, "socket4://") {
		return domain.ProxyNodeTypeSocket
	} else if strings.HasPrefix(u, "socket5://") {
		return domain.ProxyNodeTypeSocket
	}

	return domain.ProxyNodeTypeNil
}
