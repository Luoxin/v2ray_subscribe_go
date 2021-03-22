package node

import (
	"crypto/sha512"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Luoxin/v2ray_subscribe_go/conf"
	"github.com/Luoxin/v2ray_subscribe_go/db"
	"github.com/Luoxin/v2ray_subscribe_go/domain"
	"github.com/Luoxin/v2ray_subscribe_go/proxy"
	"github.com/Luoxin/v2ray_subscribe_go/utils"
)

func AddNode(nodeUrl string) (bool, error) {
	return AddNodeWithDetail(nodeUrl, 0, 0)
}

func AddNodeWithDetail(ru string, crawlerId uint64, checkInterval uint32) (bool, error) {
	if checkInterval == 0 {
		checkInterval = conf.Config.ProxyCheck.CheckInterval
	}

	proxyNodeType := utils.GetProxyNodeType(ru)

	node := &domain.ProxyNode{
		CrawlId: crawlerId,

		LastCrawlerAt: utils.Now(),
		CheckInterval: checkInterval,
		ProxyNodeType: proxyNodeType,
	}

	proxyNode, err := proxy.ParseProxy(ru)
	if err != nil {
		return false, err
	}

	proxyNode.SetCountry("")
	proxyNode.SetName("proxy")

	node.Url = proxyNode.Link()

	node.UrlFeature = fmt.Sprintf("%x", sha512.Sum512([]byte(node.Url)))

	var isNew bool
	var oldNode domain.ProxyNode
	err = db.Db.Where("url_feature = ?", node.UrlFeature).First(&oldNode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建
			log.Infof("add new proxy node: %v", node.Url)
			node.CreatedAt = utils.Now()
			err = db.Db.Create(&node).Error
			if err != nil {
				log.Errorf("err:%v", err)
				return false, err
			}
			isNew = true
		} else {
			log.Errorf("err:%v", err)
			return false, err
		}
	} else {
		// 更新
		log.Infof("update proxy node: %v", node.Url)

		node.Id = oldNode.Id
		node.CheckInterval = oldNode.CheckInterval

		node.ProxyNetworkDelay = oldNode.ProxyNetworkDelay
		node.ProxySpeed = oldNode.ProxySpeed
		node.NextCheckAt = oldNode.NextCheckAt

		if oldNode.DeathCount > 20 {
			node.DeathCount = node.DeathCount - 10
		} else {
			node.DeathCount = oldNode.DeathCount
		}
		node.AvailableCount = oldNode.AvailableCount

		err = db.Db.Save(node).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return false, err
		}
	}

	return isNew, nil
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

func GetNode4Tohru(limit int) (string, error) {
	var nodeList domain.ProxyNodeList
	err := db.Db.Select("url").Order("available_count DESC").Limit(limit).Find(&nodeList).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	var urlList []string
	nodeList.Each(func(proxyNode *domain.ProxyNode) {
		urlList = append(urlList, proxyNode.Url)
	})

	str, err := conf.Ecc.ECCEncrypt(urlList)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}

	return str, nil
}
