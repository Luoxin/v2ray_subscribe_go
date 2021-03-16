package main

import (
	log "github.com/sirupsen/logrus"

	"subscribe"
	"subscribe/tohru"
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
}
