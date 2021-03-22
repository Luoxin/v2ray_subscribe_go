package main

import (
	"github.com/luoxin/v2ray_subscribe_go/faker"
	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"
)

func main() {
	opt := &nic.H{
		AllowRedirect: true,
	}

	opt.Headers = nic.KV{
		"User-Agent": faker.New().UserAgent(),
	}

	opt.Proxy = "http://127.0.0.1:7890"

	resp, err := nic.Get("https://www.google.com", opt)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	log.Info(resp.StatusCode)
}
