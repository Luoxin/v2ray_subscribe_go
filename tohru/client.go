package tohru

import (
	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	"subscribe/conf"
)

const Hello = "B3vUNO|I,|\"FAco9b<fIPj0K:r,Zsj\"?KFOA}.z1N&LZOP1GYq"

type tohru struct {
}

func newTohru() *tohru {
	return &tohru{}
}

var Tohru = newTohru()

func (p *tohru) Init() error {

	return nil
}

func (p *tohru) DoRequest() error {
	resp, err := nic.Post(conf.Config.Base.KobayashiSanAddr+"/tohru/CheckUsable", nic.H{
		JSON: nic.KV{
			"nic": "nic",
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info(resp.Text)

	return nil
}
