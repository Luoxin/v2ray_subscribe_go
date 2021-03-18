package node

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"subscribe/conf"
	"subscribe/db"
	"subscribe/domain"
	"subscribe/proxy"
	"subscribe/utils"
)

func AddNode(nodeUrl string) error {
	err := AddNodeWithDetail(nodeUrl, 0, 0)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func AddNodeWithDetail(ru string, crawlerId uint64, checkInterval uint32) error {
	if checkInterval == 0 {
		checkInterval = conf.Config.ProxyCheck.CheckInterval
	}

	proxyNodeType := utils.GetProxyNodeType(ru)

	node := &domain.ProxyNode{
		CrawlId: crawlerId,

		LastCrawlerAt: utils.Now(),
		CheckInterval: checkInterval,
		ProxyNodeType: uint32(proxyNodeType),
	}

	proxyNode, err := proxy.ParseProxy(ru)
	if err != nil {
		return err
	}

	proxyNode.SetCountry("")
	proxyNode.SetName("proxy")

	node.Url = proxyNode.Link()

	var oldNode domain.ProxyNode
	err = db.Db.Where("url = ?", node.Url).First(&oldNode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建
			log.Infof("add new proxy node: %v", node.Url)
			node.CreatedAt = utils.Now()
			err = db.Db.Create(&node).Error
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		} else {
			log.Errorf("err:%v", err)
			return err
		}
	} else {
		// 更新
		log.Infof("update proxy node: %v", node.Url)

		node.Id = oldNode.Id
		node.CheckInterval = oldNode.CheckInterval

		node.ProxyNetworkDelay = oldNode.ProxyNetworkDelay
		node.ProxySpeed = oldNode.ProxySpeed
		node.NextCheckAt = oldNode.NextCheckAt

		if oldNode.DeathCount > 10 {
			node.DeathCount = 10
		} else {
			node.AvailableCount = oldNode.AvailableCount
		}

		err = db.Db.Save(node).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}

func GetUsableNodeList(quantity int) (domain.ProxyNodeList, error) {
	query := db.Db.Where("is_close = ?", false).
		Where("proxy_speed > 0 ").
		// Where("proxy_node_type = 1").
		Where("available_count > 0 ").
		Where("proxy_network_delay >= 0").
		//Where("death_count < ?", 10).
		// Order("proxy_node_type").
		Order("available_count DESC").
		Order("proxy_speed DESC").
		Order("proxy_network_delay").
		Order("death_count").
		Order("last_crawler_at DESC")

	if quantity >= 0 {
		query.Limit(quantity)
	}

	var nodes domain.ProxyNodeList
	err := query.Find(&nodes).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return nodes, err
}
