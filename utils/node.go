package utils

import (
	"strings"
	"subsrcibe/subscription"
)

func GetProxyNodeType(u string) subscription.ProxyNodeType {
	if strings.HasPrefix(u, "vmess") {
		return subscription.ProxyNodeType_ProxyNodeTypeVmess
	} else if strings.HasPrefix(u, "trojan") {
		return subscription.ProxyNodeType_ProxyNodeTypeTrojan
	} else if strings.HasPrefix(u, "vless") {
		return subscription.ProxyNodeType_ProxyNodeTypeVless
	}

	return subscription.ProxyNodeType_ProxyNodeTypeNil
}
