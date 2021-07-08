package node

import (
	"crypto/sha512"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/proxy"
	"github.com/Luoxin/Eutamias/utils"
)

func AddNodeWithUrl(nodeUrl string) (bool, error) {
	return AddNodeWithUrlDetail(nodeUrl, 0, 0, domain.UseTypeGFW)
}

func AddNodeWithUrlDetail(ru string, crawlerId uint64, checkInterval uint32, useType domain.UseType) (bool, error) {
	if checkInterval == 0 {
		checkInterval = conf.Config.ProxyCheck.CheckInterval
	}

	proxyNodeType := utils.GetProxyNodeType(ru)

	node := &domain.ProxyNode{
		CrawlId: crawlerId,

		LastCrawlerAt: utils.Now(),
		CheckInterval: checkInterval,
		ProxyNodeType: proxyNodeType,
		UseType:       useType,
	}

	proxyNode, err := proxy.ParseProxy(ru)
	if err != nil {
		return false, err
	}

	proxyNode.SetCountry("")
	proxyNode.SetName("proxy")

	node.Url = proxyNode.Link()

	node.UrlFeature = fmt.Sprintf("%x", sha512.Sum512([]byte(node.Url)))

	return UpdateNode(node)
}

func UpdateNode(node *domain.ProxyNode) (bool, error) {
	var isNew bool
	var oldNode domain.ProxyNode
	err := db.Db.Where("url_feature = ?", node.UrlFeature).
		First(&oldNode).Error
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
		node.IsClose = oldNode.IsClose

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

func GetUsableNodeList(quantity int, mustUsable bool, useType domain.UseType) (domain.ProxyNodeList, error) {
	query := db.Db.Where("is_close = ?", false).
		Order("available_count DESC").
		Order("proxy_speed DESC").
		Order("proxy_network_delay").
		Order("death_count").
		Order("last_crawler_at DESC")

	if useType != domain.UseTypeUseNil {
		query = query.Where("use_type = ?", uint32(useType))
	}

	if quantity >= 0 {
		query = query.Limit(quantity)
	}

	if mustUsable {
		query = query.Where("proxy_speed > 0 ").
			Where("available_count > 0 ").
			Where("proxy_network_delay >= 0").
			Where("death_count < 3")
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
	err := db.Db.Select("url").
		Where("use_type = 1").
		Order("available_count DESC").
		Limit(limit).
		Find(&nodeList).Error
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
