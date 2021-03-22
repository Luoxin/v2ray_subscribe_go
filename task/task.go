package task

import (
	"sync"
	"time"

	"github.com/roylee0704/gron"
	"github.com/roylee0704/gron/xtime"
	log "github.com/sirupsen/logrus"
	"github.com/whiteshtef/clockwork"

	"github.com/luoxin/subscribe/conf"
	"github.com/luoxin/subscribe/db"
	"github.com/luoxin/subscribe/proxycheck"
)

func InitWorker() error {
	finishC := make(chan bool, 1)
	err := InitProxy(finishC)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	sched := clockwork.NewScheduler()

	select {
	case <-finishC:
	case <-time.After(time.Minute * 10):
		log.Warn("proxy start timeout")
	}

	err = InitTohru()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	c := gron.New()

	var w sync.WaitGroup
	if conf.Config.Crawler.Enable {
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

	if conf.Config.ProxyCheck.Enable {
		log.Info("register proxy check")
		proxyCheck := proxycheck.NewProxyCheck()
		err := proxyCheck.Init()
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
			err := db.Db.
				Where("death_count > ?", 20).
				Updates(map[string]interface{}{
					"death_count": "death_count - 10",
				}).Error
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}
		})

	} else {
		log.Warnf("proxy check not start")
	}

	c.Start()

	go func() {
		sched.Run()
	}()

	return nil
}
