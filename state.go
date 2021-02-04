package main

import (
	"subsrcibe/proxycheck"
	"sync"

	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
	log "github.com/sirupsen/logrus"
	"github.com/whiteshtef/clockwork"
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
	c := gron.New()

	sched := clockwork.NewScheduler()

	var w sync.WaitGroup
	if !s.Config.DisableCrawl {
		log.Info("register crawler")

		w.Add(1)
		go func() {
			defer w.Done()

			err := crawler()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}()

		c.AddFunc(gron.Every(xtime.Minute*10), func() {
			err := crawler()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		})
	} else {
		log.Warnf("crawler not start")
	}

	if !s.Config.DisableCheckAlive {
		log.Info("register proxy check")
		proxyCheck := proxycheck.NewProxyCheck()
		err := proxyCheck.Start()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		go func() {
			w.Wait()
			w.Add(1)
			defer w.Done()

			err := checkProxyNode(proxyCheck)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}()

		c.AddFunc(gron.Every(xtime.Minute*10), func() {
			err := checkProxyNode(proxyCheck)
			if err != nil {
				log.Errorf("err:%v", err)
			}
		})

		sched.Schedule().EverySingle().Friday().At("00:00").Do(func() {
			err := s.Db.
				Where("death_count > ?", 40).
				Where("available_count > ?", 0).
				Updates(map[string]interface{}{
					"death_count": 0,

					"available_count": 0,
				}).Error
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		})

	} else {
		log.Warnf("proxy chec not start")
	}

	c.Start()

	go func() {
		sched.Run()
	}()

	return nil
}
