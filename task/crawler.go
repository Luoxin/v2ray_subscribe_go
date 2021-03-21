package task

import (
	"errors"
	"net/http"
	"time"

	"github.com/Luoxin/faker"

	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	conf2 "subscribe/conf"
	"subscribe/db"
	"subscribe/domain"
	"subscribe/node"
	"subscribe/parser"
	"subscribe/utils"
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
	defer log.Infof("crawler used %v", time.Now().Sub(t))

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
					Timeout: 60,

					SkipVerifyTLS: true,
					AllowRedirect: true,

					DisableKeepAlives: true,
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
						domain.CrawlTypeXpath,
						domain.CrawlTypeClashProxies:
						p = parser.NewFuzzyMatchingParser()
					default:
						return errors.New("nonsupport parser type")
					}

					p.ParserText(resp.Text).Each(func(nodeUrl string) {
						_, err = node.AddNodeWithDetail(nodeUrl, conf.Id, conf.Interval)
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
