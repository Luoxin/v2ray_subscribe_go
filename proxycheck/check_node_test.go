package proxycheck

import (
	"encoding/json"
	"github.com/Dreamacro/clash/adapters/outbound"
	log "github.com/sirupsen/logrus"
	"subsrcibe/utils"
	"testing"
)

func TestProxyCheck(t *testing.T) {
	proxyConfig := utils.ParseProxy("vmess://ew0KICAidiI6ICIyIiwNCiAgInBzIjogIjQ3Ljc1LjQ5LjMiLA0KICAiYWRkIjogIjQ3Ljc1LjQ5LjMiLA0KICAicG9ydCI6ICIzNjY0NCIsDQogICJpZCI6ICI3ZjE4OWFkNi0xNjBmLTRhMWYtYTg4MC1lYTdmODc2YWVhZmIiLA0KICAiYWlkIjogIjIzMyIsDQogICJuZXQiOiAid3MiLA0KICAidHlwZSI6ICJub25lIiwNCiAgImhvc3QiOiAiIiwNCiAgInBhdGgiOiAiIiwNCiAgInRscyI6ICIiDQp9")

	j, err := json.Marshal(proxyConfig)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	var proxyItem map[string]interface{}
	err = json.Unmarshal(j, &proxyItem)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	proxyItem["name"] = "test proxy"

	proxy, err := outbound.ParseProxy(proxyItem)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	//log.Info(URLTest(proxy, "https://www.google.com"))
	log.Info(URLTest(proxy, "http://cachefly.cachefly.net/10mb.test"))
	//log.Info(URLTest(proxy, "https://www.gstatic.com/generate_204"))
}
