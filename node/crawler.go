package node

import (
	"crypto/sha512"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/luoxin/v2ray_subscribe_go/subscribe/conf"
	"github.com/luoxin/v2ray_subscribe_go/subscribe/db"
	"github.com/luoxin/v2ray_subscribe_go/subscribe/domain"
	"github.com/luoxin/v2ray_subscribe_go/subscribe/utils"
)

func AddCrawlerNode(crawlerUrl string, crawlerType domain.CrawlType, rule *domain.CrawlerConf_Rule) error {
	if crawlerUrl == "" {
		return nil
	}

	node := &domain.CrawlerConf{
		CrawlerFeature: "",
		CrawlUrl:       crawlerUrl,
		CrawlType:      crawlerType,
		Rule:           rule,
		Interval:       conf.Config.Crawler.CrawlerInterval,
	}

	node.CrawlerFeature = fmt.Sprintf("%x", sha512.Sum512([]byte(node.CrawlUrl)))

	var oldNode domain.CrawlerConf
	err := db.Db.Where("crawler_feature = ?", node.CrawlerFeature).First(&oldNode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建
			log.Infof("add new proxy node: %v", node.CrawlUrl)
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
		log.Infof("update proxy node: %v", node.CrawlUrl)

		node.Id = oldNode.Id
		node.NextAt = oldNode.NextAt
		node.Interval = oldNode.Interval

		err = db.Db.Save(node).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}
