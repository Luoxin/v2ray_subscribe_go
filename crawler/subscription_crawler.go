package crawler

import (
	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"subsrcibe/subscription"
)

type SubscriptionCrawler struct {
	conf    *subscription.CrawlerConf
	proxies string
}

func (p *SubscriptionCrawler) Download() error {
	conf := p.conf

	opt := &nic.H{
		Timeout: 60,

		DisableKeepAlives:  true,
		DisableCompression: true,
		SkipVerifyTLS:      true,
	}

	rule := conf.Rule
	if rule != nil {
		if rule.UseProxy {
			opt.Proxy = p.proxies
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
}

func (p *SubscriptionCrawler) Parse() {

}
