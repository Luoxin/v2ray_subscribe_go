package main

import (
	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type State struct {
	Config *Config
	Db     *gorm.DB
}

var s *State

func initState() error {
	s = &State{
		Config: &Config{},
	}

	err := initConfig()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = initDb()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = worker()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func worker() error {
	go func() {
		c := gron.New()

		if !s.Config.DisableCrawl {
			log.Info("register crawler")
			c.AddFunc(gron.Every(xtime.Minute*10), func() {
				err := crawler()
				if err != nil {
					log.Errorf("err:%v", err)
				}
			})
		}

		if !s.Config.DisableCheckAlive {
			log.Info("register proxy check")
			c.AddFunc(gron.Every(xtime.Minute*10), func() {
				err := checkProxyNode()
				if err != nil {
					log.Errorf("err:%v", err)
				}
			})
		}

		c.Start()
	}()

	return nil
}
