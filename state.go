package main

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
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
	if !(s.Config.DisableCrawl && s.Config.DisableCheckAlive) {

		go func() {
			var err error

			err = crawler()
			if err != nil {
				log.Errorf("err:%v", err)
			}

			err = checkNode()
			if err != nil {
				log.Errorf("err:%v", err)
			}

			crawlerTicker := time.NewTimer(time.Minute * 5)
			checkTicker := time.NewTimer(time.Minute * 5)

			for {
				select {
				case <-crawlerTicker.C:
					err = crawler()
					if err != nil {
						log.Errorf("err:%v", err)
					}
					crawlerTicker.Reset(time.Minute * 5)
				case <-checkTicker.C:
					err = checkNode()
					if err != nil {
						log.Errorf("err:%v", err)
					}

					checkTicker.Reset(time.Minute * 5)
				}
			}

		}()

	} else {
		err := initCrawler()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = initCheckProxy()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return nil
}
