package eutamias

import (
	"github.com/Luoxin/Eutamias/cache"
	"github.com/Luoxin/Eutamias/dns"
	"github.com/Luoxin/Eutamias/geolite"
	"github.com/Luoxin/Eutamias/proxies"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/webservice"
	"github.com/Luoxin/Eutamias/worker"
)

func Init(configFilePatch string) error {
	err := conf.InitConfig(configFilePatch)
	if err != nil {
		log.Fatalf("init config err:%v", err)
		return err
	}

	log.Info("init conf success")

	err = cache.InitCache()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info("init cache success")

	err = proxies.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info("init clash tpl watch success")

	err = dns.InitDnsClient()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info("init dns client success")

	err = dns.InitDnsService()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info("init dns service success")

	err = geolite.InitGeoLite()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info("init geolite success")

	err = db.InitDb()
	if err != nil {
		log.Fatalf("init db err:%v", err)
		return err
	}

	log.Info("init db success")

	err = worker.InitWorker()
	if err != nil {
		log.Fatalf("init work err:%v", err)
		return err
	}

	log.Info("init worker success")

	err = webservice.InitHttpService()
	if err != nil {
		log.Fatalf("init http service err:%v", err)
		return err
	}

	log.Info("init http service success")

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
