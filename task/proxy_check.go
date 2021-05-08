package task

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/proxycheck"
	"github.com/Luoxin/Eutamias/utils"
)

func checkProxyNode() error {

	check := func(useType domain.UseType) error {
		startAt := time.Now()
		defer log.Infof("check proxy used %v for %v", time.Since(startAt), domain.UseTypeMap[useType])

		var nodeList domain.ProxyNodeList
		err := db.Db.Where("is_close = ?", false).
			Where("next_check_at < ?", utils.Now()).
			Where("death_count < ?", 50).
			Where("use_type = ?", uint32(useType)).
			Order("next_check_at").
			Find(&nodeList).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		log.Infof("check proxy for %v node for %v", len(nodeList), domain.UseTypeMap[useType])

		check := proxycheck.NewProxyCheck()
		err = check.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		switch useType {
		case domain.UseTypeNetEase:
			check.SetCheckUrl("http://interface.music.163.com")
		}

		nodeList.Each(func(node *domain.ProxyNode) {
			err = check.AddWithLink(node.Url, func(result proxycheck.Result) error {
				um := map[string]interface{}{}

				if result.Err != nil {
					node.DeathCount++
					if node.DeathCount > 10 {
						node.ProxySpeed = -1
						node.ProxyNetworkDelay = -1
						node.AvailableCount = 0

						um["proxy_speed"] = -1
						um["proxy_network_delay"] = -1
						um["available_count"] = 0
					}

					um["death_count"] = node.DeathCount
				} else {
					node.DeathCount = 0
					node.AvailableCount++
					node.ProxyNetworkDelay = result.Delay
					node.ProxySpeed = result.Speed

					um["available_count"] = node.AvailableCount
					um["death_count"] = 0

					um["proxy_speed"] = result.Speed
					um["proxy_network_delay"] = result.Delay
				}

				log.Infof("check proxy %+v: speed:%.3f Kb/s, delay:%v ms, available %d, death %d",
					utils.ShortStr(node.UrlFeature, 12), node.ProxySpeed, node.ProxyNetworkDelay, node.AvailableCount, node.DeathCount)

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

	err := check(domain.UseTypeUseNil)
	if err != nil {
		log.Errorf("err:%v", err)
	}

	err = check(domain.UseTypeGFW)
	if err != nil {
		log.Errorf("err:%v", err)
	}

	err = check(domain.UseTypeNetEase)
	if err != nil {
		log.Errorf("err:%v", err)
	}

	return nil
}
