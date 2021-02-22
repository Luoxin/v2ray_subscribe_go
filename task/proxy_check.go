package task

import (
	log "github.com/sirupsen/logrus"
	"subscribe/db"
	"subscribe/domain"
	"subscribe/proxycheck"
	"subscribe/utils"
	"time"
)

func checkProxyNode(check *proxycheck.ProxyCheck) error {
	t := time.Now()
	log.Info("start check proxy...")
	defer log.Infof("check proxy used %v", time.Now().Sub(t))

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

	if len(nodeList) == 0 {
		log.Warnf("not found nodes need check")
		return nil
	}

	nodeList.Each(func(node *domain.ProxyNode) {
		if node.NodeDetail == nil {
			// TODO 移除节点
			return
		}

		err = check.AddWithLink(node.NodeDetail.Buf, func(result proxycheck.Result) error {
			um := map[string]interface{}{}

			if result.Err != nil {
				//log.Info(reflect.TypeOf(result.Err))
				//log.Errorf("err:%v", result.Err)

				node.DeathCount++
				if node.DeathCount > 10 {
					node.ProxySpeed = -1
					node.ProxyNetworkDelay = -1
				}
			} else {
				node.DeathCount = 0
				node.AvailableCount++

				um["available_count"] = node.AvailableCount
			}

			um["death_count"] = node.DeathCount
			um["proxy_speed"] = result.Speed
			um["proxy_network_delay"] = result.Delay

			log.Infof("check proxy %+v: speed:%v Mb/s, delay:%v ms,available %d, death %d",
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
