package crawler

import "github.com/Luoxin/Eutamias/domain"

type Downloader interface {
	Download(method string, urlStr string, reqBody interface{}, rule domain.CrawlerConf_Rule) (string, error)
}
