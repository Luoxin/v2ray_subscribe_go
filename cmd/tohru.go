package main

import (
	log "github.com/sirupsen/logrus"

	"subscribe"
	"subscribe/conf"
	"subscribe/db"
	"subscribe/tohru"
	"subscribe/webservice"
)

func main() {
	err := conf.InitConfig()
	if err != nil {
		log.Errorf("err:%v", err)
	}

	err = db.InitDb(conf.Config.Db.Addr)
	if err != nil {
		log.Errorf("err:%v", err)
	}

	return

	err = conf.InitConfig()
	if err != nil {
		log.Errorf("err:%v", err)
	}

	err = webservice.InitHttpService()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	return

	err = subscribe.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = tohru.Tohru.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
