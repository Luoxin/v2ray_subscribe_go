package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"subsrcibe/parser"

	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"subsrcibe/domain"
	"subsrcibe/utils"
)

func crawler() error {
	log.Infof("start crawler......")

	var crawlerList []*domain.CrawlerConf
	err := s.Db.Where("is_closed = ?", false).
		Where("next_at < ?", utils.Now()).
		Order("next_at").
		//Limit(1).
		Find(&crawlerList).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	CrawlerConfList(crawlerList).
		Each(func(conf *domain.CrawlerConf) {
			if conf.CrawlUrl == "" {
				log.Errorf("crawler url empty: %d", conf.Id)
				return
			}

			log.Infof("wail crawler for %+v", conf.CrawlUrl)
			defer log.Infof("crawler finish,%v next exec at %v", conf.CrawlUrl, conf.NextAt)

			err := func() error {
				opt := &nic.H{
					Timeout: 60,

					SkipVerifyTLS: true,
					AllowRedirect: true,

					DisableKeepAlives: true,
				}

				rule := conf.Rule
				if rule != nil {
					if rule.UseProxy {
						opt.Proxy = s.Config.Proxies
					}
				}

				resp, err := nic.Get(conf.CrawlUrl, opt)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}

				switch resp.StatusCode {
				case http.StatusOK:
					var p parser.Parser

					switch domain.CrawlType(conf.CrawlType) {
					case domain.CrawlType_CrawlTypeSubscription, domain.CrawlType_CrawlTypeFuzzyMatching, domain.CrawlType_CrawlTypeXpath:
						p = parser.NewFuzzyMatchingParser()
					//case domain.CrawlType_CrawlTypeXpath:
					default:
						return errors.New("nonsupport parser type")
					}

					p.ParserText(resp.Text).Each(func(nodeUrl string) {
						err = addNode(nodeUrl, conf.Id, conf.Interval)
						if err != nil {
							log.Errorf("err:%v", err)
							return
						}

					})

				case http.StatusMovedPermanently, http.StatusFound:
					// 重定向了
					u, err := resp.Location()
					if err != nil {
						log.Errorf("err:%v", err)
						return err
					}

					log.Warnf("%v redirect to %v", conf.CrawlUrl, u)
					conf.CrawlUrl = u.String()
					return nil

				case http.StatusNonAuthoritativeInfo:
					// 不可信的信息
					return nil

				default:
					log.Warnf("not support status code %v", resp.StatusCode)
					return nil
				}

				return nil
			}()
			if err != nil {
				log.Errorf("err:%v", err)
				conf.NextAt = conf.Interval*2 + utils.Now()
			}

			if conf.NextAt < utils.Now() {
				if conf.Interval == 0 {
					conf.Interval = s.Config.CrawlerInterval
				}

				conf.NextAt = conf.Interval + utils.Now()
			}

			// 保存数据
			err = s.Db.Omit("created_at").Save(conf).Error
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		})

	return nil
}

func addNode(ru string, crawlerId uint64, checkInterval uint32) error {
	if checkInterval == 0 {
		checkInterval = s.Config.CheckInterval
	}

	proxyNodeType := utils.GetProxyNodeType(ru)

	node := &domain.ProxyNode{
		NodeDetail: &domain.ProxyNode_NodeDetail{
			Buf: ru,
		},
		CrawlId: crawlerId,

		LastCrawlerAt: utils.Now(),
		CheckInterval: checkInterval,
		ProxyNodeType: uint32(proxyNodeType),
	}

	nodeInterface := utils.ParseProxy(ru)

	switch proxyNodeType {
	case domain.ProxyNodeType_ProxyNodeTypeVmess:
		nodeItem := nodeInterface.(utils.ClashVmess)

		node.Url = fmt.Sprintf("%v:%v/", nodeItem.Server, nodeItem.Port)
		if nodeItem.Network == "ws" {
			node.Url += strings.TrimPrefix(nodeItem.WSPATH, "/")
		}

	case domain.ProxyNodeType_ProxyNodeTypeTrojan:
		nodeItem := nodeInterface.(utils.Trojan)
		node.Url = fmt.Sprintf("%v:%v/", nodeItem.Server, nodeItem.Port)

	//case domain.ProxyNodeType_ProxyNodeTypeVless:
	case domain.ProxyNodeType_ProxyNodeTypeSS:
		nodeItem := nodeInterface.(utils.ClashSS)
		node.Url = fmt.Sprintf("%v:%v/", nodeItem.Server, nodeItem.Port)

	case domain.ProxyNodeType_ProxyNodeTypeSSR:
		nodeItem := nodeInterface.(utils.ClashRSSR)
		node.Url = fmt.Sprintf("%v:%v/", nodeItem.Server, nodeItem.Port)

	default:
		return ErrInvalidArg
	}

	var oldNode domain.ProxyNode
	err := s.Db.Where("url = ?", node.Url).First(&oldNode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建
			log.Infof("add new proxy node: %v", node.Url)
			node.CreatedAt = utils.Now()
			err = s.Db.Create(&node).Error
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

		err = s.Db.Save(node).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}
