package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/eddieivan01/nic"
	"github.com/elliotchance/pie/pie"
	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"subsrcibe/subscribe"
	"subsrcibe/utils"
	"time"
)

func initCrawler() error {
	if s.Config.DisableCrawl {
		log.Warnf("crawler disable")
		return nil
	}

	go func() {
		log.Info("start crawler work...")
		for {
			err := crawler()
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}

			time.Sleep(time.Minute * 10)
		}

	}()

	return nil
}

func crawler() error {
	var crawlerList []*subscribe.CrawlerConf
	err := s.Db.Where("is_closed = ?", false).
		Where("next_at < ?", utils.Now()).
		//Limit(1).
		Find(&crawlerList).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	CrawlerConfList(crawlerList).Each(func(conf *subscribe.CrawlerConf) {
		if conf.CrawlUrl == "" {
			log.Errorf("crawler url empty: %d", conf.Id)
			return
		}

		log.Infof("wail crawler for %+v", conf.CrawlUrl)
		defer log.Infof("crawler finish, next exec at %v", conf.NextAt)

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
				switch subscribe.CrawlType(conf.CrawlType) {
				case subscribe.CrawlType_CrawlTypeSubscription:
					//log.Infof("get node info %v", resp.Text)
					err = addNodesByBase64(conf, resp.Text)
					if err != nil {
						log.Errorf("err:%v", err)
						return err
					}

				case subscribe.CrawlType_CrawlTypeXpath:
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
				conf.NextAt += conf.Interval*10 + utils.Now()
				return nil

			default:
				log.Warnf("not support status code %v", resp.StatusCode)
				return nil
			}

			return nil
		}()
		if err != nil {
			log.Errorf("err:%v", err)
		}

		// 保存数据
		err = s.Db.Save(conf).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})

	return nil
}

func addNodesByBase64(crawlerConf *subscribe.CrawlerConf, bs string) error {
	if bs == "" {
		log.Warnf("nodes empty")
		return nil
	}

	decode := func(s string) (string, error) {
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

	str, err := decode(bs)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	urlList := strings.Split(str, "\n")

	if crawlerConf == nil {
		crawlerConf = &subscribe.CrawlerConf{
			Interval: s.Config.CheckInterval,
		}
	}

	pie.Strings(urlList).Each(func(ru string) {
		node := &subscribe.ProxyNode{
			NodeDetail: &subscribe.ProxyNode_NodeDetail{
				Buf: ru,
			},
			CrawlId: crawlerConf.Id,

			CheckInterval: crawlerConf.Interval,
		}

		if !strings.HasPrefix(ru, "vmess") {
			return
		}

		node.ProxyNodeType = uint32(subscribe.ProxyNodeType_ProxyNodeTypeVmess)

		d, err := decode(strings.TrimPrefix(ru, "vmess://"))
		if err != nil {
			log.Errorf("err:%v", err)
			//log.Errorf("ru %v", ru)
			return
		}

		var vmessNode subscribe.ProxyNode_VmessNode
		err = jsonpb.Unmarshal(bytes.NewBufferString(d), &vmessNode)
		if err != nil {
			log.Errorf("err:%v", err)
			//log.Errorf("d %v", d)
			return
		}

		node.NodeDetail.VmessNode = &vmessNode

		host := vmessNode.Host
		if host == "" {
			host = vmessNode.Add
		}

		node.Url = fmt.Sprintf("%v:%v/%v", host, vmessNode.Port, strings.TrimPrefix(vmessNode.Path, "/"))

		var oldNode subscribe.ProxyNode
		err = s.Db.Where("url = ?", node.Url).First(&oldNode).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// 创建
				log.Infof("add new proxy node: %v", node.Url)

				err = s.Db.Create(&node).Error
				if err != nil {
					log.Errorf("err:%v", err)
					return
				}
			} else {
				log.Errorf("err:%v", err)
				return
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
				return
			}
		}

	})

	return nil
}
