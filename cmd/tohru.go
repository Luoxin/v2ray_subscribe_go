package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/luoxin/subscribe"
	"github.com/luoxin/subscribe/tohru"
)

func main() {
	err := subscribe.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = tohru.Tohru.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	// err = tohru.Tohru.Registration(conf.Config.Base.TohruKey, conf.Config.Base.TohruPassword)
	// if err != nil {
	// 	log.Errorf("err:%v", err)
	// 	return
	// }

	err = tohru.Tohru.Start()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
