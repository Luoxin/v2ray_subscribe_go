package tohru

import (
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"subscribe/conf"
	"subscribe/version"
)

const Hello = "B3vUNO|I,|\"FAco9b<fIPj0K:r,Zsj\"?KFOA}.z1N&LZOP1GYq"

type tohru struct {
	client *resty.Client
}

func newTohru() *tohru {
	return &tohru{}
}

var Tohru = newTohru()

func (p *tohru) Init() error {
	p.client = resty.New().
		SetTimeout(time.Second * 5).
		SetRetryMaxWaitTime(time.Second * 5).
		SetRetryWaitTime(time.Second)

	if conf.Config.IsTohru() {
		err := p.CheckUsable()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}

func (p *tohru) CheckUsable() error {
	hello, err := conf.Ecc.ECCEncrypt(Hello)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var rsp CheckUsableRsp
	err = p.DoRequest("/tohru/CheckUsable", CheckUsableReq{
		Version: version.Version,
		Hello:   hello,
	}, &rsp)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (p *tohru) DoRequest(path string, req, rsp interface{}) error {
	resp, err := p.client.R().
		SetBody(req).
		SetResult(&rsp).
		Post(conf.Config.Base.KobayashiSanAddr + path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info(resp.Result())

	// resp, err := nic.Post(, nic.H{
	// 	Headers:            nil,
	// 	Cookies:            nil,
	// 	Auth:               nil,
	// 	Proxy:              "",
	// 	JSON:               nil,
	// 	Files:              nil,
	// 	AllowRedirect:      true,
	// 	Timeout:            5,
	// 	Chunked:            false,
	// 	DisableKeepAlives:  false,
	// 	DisableCompression: false,
	// 	SkipVerifyTLS:      false,
	// })

	return nil
}
