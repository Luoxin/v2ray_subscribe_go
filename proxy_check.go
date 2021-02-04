package main

import (
	log "github.com/sirupsen/logrus"
	"subsrcibe/domain"
	"subsrcibe/proxycheck"
	"subsrcibe/utils"
)

func checkProxyNode(check *proxycheck.ProxyCheck) error {
	var nodeList domain.ProxyNodeList
	err := s.Db.Where("is_close = ?", false).
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

		err = check.Add(node.NodeDetail.Buf, func(err error, delay, speed float64) error {
			e := func() error {
				if err != nil {
					return err
				}

				node.AvailableCount++
				return nil
			}()
			if e != nil {
				log.Errorf("err:%v", e)
				node.DeathCount++

				if node.DeathCount > 10 {
					node.ProxySpeed = -1
					node.ProxyNetworkDelay = -1
				}
			} else {
				node.DeathCount = 0
				node.AvailableCount++
				log.Infof("check proxy %+v: speed:%v, delay:%v", node.Url, node.ProxySpeed, node.ProxyNetworkDelay)
			}

			node.NextCheckAt = node.CheckInterval + utils.Now()

			err = s.Db.Omit("node_detail", "url", "proxy_node_type").Save(node).Error
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

	return nil
}
