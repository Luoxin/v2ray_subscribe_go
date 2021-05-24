package utils

import "net"

func IsIp(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsIpV4(ip string) bool {
	IP := net.ParseIP(ip)
	if IP == nil {
		return false
	}

	return IP.To4() != nil
}

func IsIpV6(ip string) bool {
	IP := net.ParseIP(ip)
	if IP == nil {
		return false
	}

	return IP.To16() != nil
}
