package download

import "github.com/luoxin/v2ray_subscribe_go/domain"

type Downloader interface {
	Download(method string, urlStr string, reqBody string, rule domain.CrawlerConf_Rule) (string, error)
}
