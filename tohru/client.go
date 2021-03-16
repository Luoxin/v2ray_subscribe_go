package tohru

import (
	"fmt"

	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	"subscribe/conf"
	"subscribe/version"
)

const Hello = "B3vUNO|I,|\"FAco9b<fIPj0K:r,Zsj\"?KFOA}.z1N&LZOP1GYq"

type tohru struct {
}

func newTohru() *tohru {
	return &tohru{}
}

var Tohru = newTohru()

func (p *tohru) Init() error {
	err := p.CheckUsable()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (p *tohru) CheckUsable() error {
	if conf.Config.IsTohru() {
		err := p.DoRequest("/tohru/CheckUsable")
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}

func (p *tohru) DoRequest(path string) error {
	hello, err := conf.Ecc.ECCEncrypt(Hello)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	fmt.Println(hello)

	resp, err := nic.Post(conf.Config.Base.KobayashiSanAddr+path, nic.H{
		JSON: nic.KV{
			"version": version.Version,
			"hello":   hello,
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info(resp.Text)

	return nil
}
