package utils

import (
	"strings"

	"github.com/Luoxin/Eutamias/domain"
)

func GetProxyNodeType(u string) domain.ProxyNodeType {
	if strings.HasPrefix(u, "vmess://") {
		return domain.ProxyNodeType_ProxyNodeTypeVmess
	} else if strings.HasPrefix(u, "trojan://") {
		return domain.ProxyNodeType_ProxyNodeTypeTrojan
	} else if strings.HasPrefix(u, "vless://") {
		return domain.ProxyNodeType_ProxyNodeTypeVless
	} else if strings.HasPrefix(u, "ssr://") {
		return domain.ProxyNodeType_ProxyNodeTypeSSR
	} else if strings.HasPrefix(u, "ss://") {
		return domain.ProxyNodeType_ProxyNodeTypeSS
	} else if strings.HasPrefix(u, "http://") {
		return domain.ProxyNodeType_ProxyNodeTypeHttp
	} else if strings.HasPrefix(u, "socket://") {
		return domain.ProxyNodeType_ProxyNodeTypeSocket
	} else if strings.HasPrefix(u, "socket4://") {
		return domain.ProxyNodeType_ProxyNodeTypeSocket
	} else if strings.HasPrefix(u, "socket5://") {
		return domain.ProxyNodeType_ProxyNodeTypeSocket
	}

	return domain.ProxyNodeType_ProxyNodeTypeNil
}
