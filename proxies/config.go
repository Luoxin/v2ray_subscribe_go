package proxies

import (
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/node"
	log "github.com/sirupsen/logrus"
)

func GenClashConfig(count int, mustUsable, needCheck bool, c ClashConfig) (string, int) {
	nodes, err := node.GetUsableNodeList(count, mustUsable, domain.UseTypeGFW)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", 0
	}

	p := NewProxies()

	nodes.Each(func(node *domain.ProxyNode) {
		p.AppendWithUrl(node.Url)
	})

	if needCheck {
		p = p.GetUsableList()
	}

	nodes, err = node.GetUsableNodeList(count, mustUsable, domain.UseTypeNetEase)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", 0
	}
	nodes.Each(func(proxyNode *domain.ProxyNode) {
		p.AppendNetEaseWithUrl(proxyNode.Url)
	})

	return p.ToClashConfig(c), p.Len()
}
