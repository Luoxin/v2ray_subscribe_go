package proxies

import (
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/node"
	log "github.com/sirupsen/logrus"
)

func GenClashConfig(count int, mustUsable, needCheck bool) (string, int) {
	nodes, err := node.GetUsableNodeList(count, mustUsable)
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

	err = db.Db.Where("is_close = ?", false).
		Where("use_type = 2").Find(&nodes).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return "", 0
	}

	nodes.Each(func(proxyNode *domain.ProxyNode) {
		p.AppendNetEaseWithUrl(proxyNode.Url)
	})

	return p.ToClashConfig(), p.Len()
}
