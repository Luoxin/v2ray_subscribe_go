package utils

import "net"

func IsIpAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}
