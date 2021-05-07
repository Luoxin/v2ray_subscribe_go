package task

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Luoxin/faker"
	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	conf2 "github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/node"
	"github.com/Luoxin/Eutamias/parser"
	"github.com/Luoxin/Eutamias/utils"
)

func crawler() error {
	t := time.Now()

	var crawlerList []*domain.CrawlerConf
	err := db.Db.Where("is_closed = ?", false).
		Where("next_at < ?", utils.Now()).
		Order("next_at").
		//Limit(1).
		Find(&crawlerList).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("crawler for %v website", len(crawlerList))
	defer log.Infof("crawler used %v", time.Since(t))

	domain.CrawlerConfList(crawlerList).
		Each(func(conf *domain.CrawlerConf) {
			if conf.CrawlUrl == "" {
				log.Errorf("crawler url empty: %d", conf.Id)
				return
			}

			log.Infof("wail crawler for %+v", conf.CrawlUrl)
			defer log.Infof("crawler finish,%v next exec at %v", conf.CrawlUrl, conf.NextAt)

			err := func() error {
				opt := &nic.H{
					AllowRedirect:     true,
					Timeout:           60,
					DisableKeepAlives: true,
					SkipVerifyTLS:     true,
				}

				opt.Headers = nic.KV{
					"User-Agent": faker.New().UserAgent(),
				}

				rule := conf.Rule
				if rule != nil {
					if rule.UseProxy {
						opt.Proxy = conf2.Config.Crawler.Proxies
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

					switch conf.CrawlType {
					case domain.CrawlTypeSubscription,
						domain.CrawlTypeFuzzyMatching,
						domain.CrawlTypeXpath:
						p = parser.NewFuzzyMatchingParser()
					case domain.CrawlTypeClashProxies:
						p = parser.NewClashParser()
					default:
						return errors.New("nonsupport parser type")
					}

					_ = p.ParserText(resp.Text).Filter(func(s string) bool {
						return strings.Contains(s, "://")
					}).Each(func(nodeUrl string) {
						_, err = node.AddNodeWithUrlDetail(nodeUrl, conf.Id, conf.Interval, conf.UseType)
						if err != nil {
							log.Errorf("link:%s, err:%v", nodeUrl, err)
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

				case http.StatusNotFound:
					log.Warnf("%v is not found,will closed", conf.CrawlUrl)
					conf.IsClosed = true
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
					conf.Interval = conf2.Config.Crawler.CrawlerInterval
				}

				conf.NextAt = conf.Interval + utils.Now()
			}

			// 保存数据
			err = db.Db.Omit("created_at").Save(conf).Error
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		})

	return nil
}
