package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/eddieivan01/nic"
	"github.com/elliotchance/pie/pie"
	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"strings"
	"subsrcibe/subscription"
	"subsrcibe/utils"
	"time"
)

func initCrawler() error {
	if s.Config.DisableCrawl {
		log.Warnf("crawler is disable")
		return nil
	}

	go func() {
		log.Info("start crawler worker......")
		for {
			err := crawler()
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}

			time.Sleep(time.Minute * 5)
		}

	}()

	return nil
}

func crawler() error {
	log.Infof("start crawler......")

	var crawlerList []*subscription.CrawlerConf
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
		Each(func(conf *subscription.CrawlerConf) {
			if conf.CrawlUrl == "" {
				log.Errorf("crawler url empty: %d", conf.Id)
				return
			}

			log.Infof("wail crawler for %+v", conf.CrawlUrl)
			defer log.Infof("crawler finish,%v next exec at %v", conf.CrawlUrl, conf.NextAt)

			err := func() error {
				opt := &nic.H{
					Proxy:   s.Config.Proxies,
					Timeout: 60,

					DisableKeepAlives:  true,
					DisableCompression: true,
					SkipVerifyTLS:      true,
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
					switch subscription.CrawlType(conf.CrawlType) {
					case subscription.CrawlType_CrawlTypeSubscription:
						//log.Infof("get node info %v", resp.Text)
						err = addNodesByBase64(conf, resp.Text)
						if err != nil {
							log.Errorf("err:%v", err)
							return err
						}

					case subscription.CrawlType_CrawlTypeXpath:
					case subscription.CrawlType_CrawlTypeFuzzyMatching:
						pie.
							Strings(regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$`).
								FindStringSubmatch(resp.Text)).
							Each(func(ru string) {
								err = addNodesByBase64(conf, resp.Text)
								if err != nil {
									log.Errorf("err:%v", err)
									return
								}
							})

						pie.
							Strings(regexp.MustCompile(`vmess://[^\s]*`).
								FindStringSubmatch(resp.Text)).
							Each(func(ru string) {
								err = addNode(ru, conf.Id, conf.Interval)
								if err != nil {
									log.Errorf("err:%v", err)
									return
								}
							})

						pie.
							Strings(regexp.MustCompile(`trojan://[^\s]*`).
								FindStringSubmatch(resp.Text)).
							Each(func(ru string) {
								err = addNode(ru, conf.Id, conf.Interval)
								if err != nil {
									log.Errorf("err:%v", err)
									return
								}
							})

						pie.
							Strings(regexp.MustCompile(`vless://[^\s]*`).
								FindStringSubmatch(resp.Text)).
							Each(func(ru string) {
								err = addNode(ru, conf.Id, conf.Interval)
								if err != nil {
									log.Errorf("err:%v", err)
									return
								}
							})
					}

					conf.NextAt += conf.Interval + utils.Now()

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
					conf.NextAt += conf.Interval*5 + utils.Now()
					return nil

				default:
					log.Warnf("not support status code %v", resp.StatusCode)
					return nil
				}

				return nil
			}()
			if err != nil {
				conf.NextAt = conf.Interval + utils.Now()
				log.Errorf("err:%v", err)
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

func addNodesByBase64(crawlerConf *subscription.CrawlerConf, bs string) error {
	if bs == "" {
		log.Warnf("nodes empty")
		return nil
	}

	str, err := base64Decode(bs)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	urlList := strings.Split(str, "\n")

	if crawlerConf == nil {
		crawlerConf = &subscription.CrawlerConf{
			Interval: s.Config.CheckInterval,
		}
	}

	pie.Strings(urlList).Each(func(ru string) {
		err = addNode(ru, crawlerConf.Id, crawlerConf.Interval)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})

	return nil
}

func base64Decode(s string) (string, error) {
	str, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		str, err = base64.StdEncoding.DecodeString(s)
		if err != nil {
			str, err = base64.RawStdEncoding.DecodeString(s)
			if err != nil {
				str, err = base64.RawURLEncoding.DecodeString(s)
				if err != nil {
					log.Errorf("decode fail for %v", s)
					return "", err
				}
			}
		}
	}
	return string(str), err
}

func addNode(ru string, crawlerId uint64, checkInterval uint32) error {
	log.Infof("will add node:%v", ru)

	if checkInterval == 0 {
		checkInterval = s.Config.CheckInterval
	}

	proxyNodeType := utils.GetProxyNodeType(ru)

	node := &subscription.ProxyNode{
		NodeDetail: &subscription.ProxyNode_NodeDetail{
			Buf: ru,
		},
		CrawlId: crawlerId,

		LastCrawlerAt: utils.Now(),
		CheckInterval: checkInterval,
		ProxyNodeType: uint32(proxyNodeType),
	}

	switch proxyNodeType {
	case subscription.ProxyNodeType_ProxyNodeTypeVmess:
		d, err := base64Decode(strings.TrimPrefix(ru, "vmess://"))
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		var vmessNode subscription.ProxyNode_VmessNode
		err = jsonpb.Unmarshal(bytes.NewBufferString(d), &vmessNode)
		if err != nil {
			log.Errorf("err:%v", err)

			m := map[string]interface{}{}
			err = json.Unmarshal([]byte(d), &m)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			for k, v := range m {
				m[strings.ToLower(k)] = v
			}

			host := m["host"]
			if host == "" {
				host = m["add"]
			}

			node.Url = fmt.Sprintf("%v:%v/%v", host, m["port"], strings.TrimPrefix(fmt.Sprintf("%v", m["path"]), "/"))

		} else {
			node.NodeDetail.VmessNode = &vmessNode
			host := vmessNode.Host
			if host == "" {
				host = vmessNode.Add
			}

			node.Url = fmt.Sprintf("%v:%v/%v", host, vmessNode.Port, strings.TrimPrefix(vmessNode.Path, "/"))
		}
	case subscription.ProxyNodeType_ProxyNodeTypeTrojan:
		node.Url = strings.Split(strings.TrimPrefix(ru, "trojan://"), "#")[0]
	default:
		return ErrInvalidArg
	}

	var oldNode subscription.ProxyNode
	err := s.Db.Where("url = ?", node.Url).First(&oldNode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建
			log.Infof("add new proxy node: %v", node.Url)

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
