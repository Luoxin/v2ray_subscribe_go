package task

import (
	log "github.com/sirupsen/logrus"
	"github.com/luoxin/v2ray_subscribe_go/db"
	"github.com/luoxin/v2ray_subscribe_go/domain"
	"github.com/luoxin/v2ray_subscribe_go/proxycheck"
	"github.com/luoxin/v2ray_subscribe_go/utils"
	"time"
)

func checkProxyNode(check *proxycheck.ProxyCheck) error {
	t := time.Now()

	var nodeList domain.ProxyNodeList
	err := db.Db.Where("is_close = ?", false).
		Where("next_check_at < ?", utils.Now()).
		Where("death_count < ?", 50).
		Order("next_check_at").
		Find(&nodeList).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("check proxy for %v node", len(nodeList))
	defer log.Infof("check proxy used %v", time.Now().Sub(t))

	nodeList.Each(func(node *domain.ProxyNode) {
		err = check.AddWithLink(node.Url, func(result proxycheck.Result) error {
			um := map[string]interface{}{}

			if result.Err != nil {
				//log.Info(reflect.TypeOf(result.Err))
				//log.Errorf("err:%v", result.Err)

				node.DeathCount++
				if node.DeathCount > 10 {
					node.ProxySpeed = -1
					node.ProxyNetworkDelay = -1
					node.AvailableCount = 0

					um["proxy_speed"] = -1
					um["proxy_network_delay"] = -1
					um["available_count"] = 0
				}
			} else {
				node.DeathCount = 0
				node.AvailableCount++
				node.ProxyNetworkDelay = result.Delay
				node.ProxySpeed = result.Speed

				um["available_count"] = node.AvailableCount

				um["proxy_speed"] = result.Speed
				um["proxy_network_delay"] = result.Delay
			}

			um["death_count"] = node.DeathCount

			log.Infof("check proxy %+v: speed:%v Kb/s, delay:%v ms, available %d, death %d",
				node.Url, node.ProxySpeed, node.ProxyNetworkDelay, node.AvailableCount, node.DeathCount)

			node.NextCheckAt = node.CheckInterval + utils.Now()

			um["next_check_at"] = node.NextCheckAt

			err = db.Db.Model(node).Updates(um).Error
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			return nil
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

	})

	check.WaitFinish()

	return nil
}
