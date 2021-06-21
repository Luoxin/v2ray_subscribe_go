package crawler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/domain"
)

type HttpDownloader struct {
	option *nic.H
}

func NewHttpDownloader() HttpDownloader {
	return HttpDownloader{
		option: &nic.H{
			Timeout:            5,
			Chunked:            false,
			DisableCompression: true,
			SkipVerifyTLS:      true,
		},
	}
}

func (h *HttpDownloader) Download(method string, urlStr string, reqBody interface{}, rule domain.CrawlerConf_Rule) (string, error) {
	opt := &nic.H{
		Timeout: 60,

		SkipVerifyTLS: true,
		AllowRedirect: true,

		DisableKeepAlives: true,
	}

	if rule.UseProxy {
		opt.Proxy = conf.Config.Crawler.Proxies
	}

	var resp *nic.Response
	var err error

	switch strings.ToUpper(method) {
	case "GET":
		resp, err = nic.Get(urlStr, opt)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}

	default:
		return "", errors.New("nonsupport method")
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.Text, nil

	case http.StatusMovedPermanently, http.StatusFound:
		// 重定向了
		u, err := resp.Location()
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}

		return u.String(), errors.New("moved permanently")

	case http.StatusNonAuthoritativeInfo:
		// 不可信的信息
		return "", errors.New("non authoritative info")

	default:
		log.Warnf("nonsupport status code %v", resp.StatusCode)
		return "", fmt.Errorf("nonsupport status code %v", resp.StatusCode)
	}
}
