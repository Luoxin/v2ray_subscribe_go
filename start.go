package eutamias

import (
	"github.com/Luoxin/Eutamias/cache"
	"github.com/Luoxin/Eutamias/dns"
	"github.com/Luoxin/Eutamias/geolite"
	"github.com/Luoxin/Eutamias/initialize"
	"github.com/Luoxin/Eutamias/proxies"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/webservice"
	"github.com/Luoxin/Eutamias/worker"
)

var ConfigFilePatch *string

var initList = []initialize.Initialize{
	&conf.Init{
		ConfigFilePatch: ConfigFilePatch,
	},
	&cache.Init{},
	&proxies.Init{},
	&db.Init{},
	&dns.InitClient{},
	&dns.InitServer{},
	&geolite.Init{},
	&worker.Init{},
	&webservice.Init{},
}

func Init(configFilePatch string) error {
	for _, item := range initList {
		needRun, err := item.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		if !needRun {

		}

		item.WaitFinish()
	}

	return nil
}

func Start() {
	c := make(chan bool)

	err := Init("")
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

	<-c
}
