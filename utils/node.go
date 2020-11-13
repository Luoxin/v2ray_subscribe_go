package utils

import (
	"strings"
	"subsrcibe/subscription"
)

func GetProxyNodeType(u string) subscription.ProxyNodeType {
	if strings.HasPrefix(u, "vmess") {
		return subscription.ProxyNodeType_ProxyNodeTypeVmess
	} else  if strings.HasPrefix(u, "trojan") {
		return subscription.ProxyNodeType_ProxyNodeTypeVmess
	}

	return subscription.ProxyNodeType_ProxyNodeTypeNil
}