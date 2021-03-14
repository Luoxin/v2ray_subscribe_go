package main

import (
	"fmt"

	"github.com/Luoxin/faker"
	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	"subscribe/utils"
)

func main() {
	s := "test"
	e, err := utils.ECCEncrypt([]byte(s))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	fmt.Println(string(e))

	d, err := utils.ECCDecrypt(e)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	fmt.Println(string(d))

	return
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
