package utils

import (
	"strings"
	"subscribe/domain"
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
	}

	return domain.ProxyNodeType_ProxyNodeTypeNil
}
