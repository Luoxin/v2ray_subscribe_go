package proxy

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestParseProxyFromClashProxy(t *testing.T) {
	proxyConfig, err := ParseProxy("vmess://ew0KICAidiI6ICIyIiwNCiAgInBzIjogIjEt5omT5bel6ZKx5bCR5bCx5Ye65p2l5Yib5Lia5ZCn77yM5L+d6K+B6IO96K6p5L2g77yM6LWU5Liq57K+5YWJ44CCIiwNCiAgImFkZCI6ICIxMDEuMzIuMTg5LjQwIiwNCiAgInBvcnQiOiAiODAiLA0KICAiaWQiOiAiMTNkN2MzN2UtOTIyZi02MWFkLTYzNmEtOWFlNDVmMmE3OGQxIiwNCiAgImFpZCI6ICIwIiwNCiAgIm5ldCI6ICJ3cyIsDQogICJ0eXBlIjogIm5vbmUiLA0KICAiaG9zdCI6ICJjZG4xNjN2My5kMi41MWp1bXAuY28iLA0KICAicGF0aCI6ICIvdXNtYXdqZHAiLA0KICAidGxzIjogIiINCn0=")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	log.Info(proxyConfig.BaseInfo().Server)
}
