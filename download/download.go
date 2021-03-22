package download

import "github.com/luoxin/subscribe/domain"

type Downloader interface {
	Download(method string, urlStr string, reqBody string, rule domain.CrawlerConf_Rule) (string, error)
}
