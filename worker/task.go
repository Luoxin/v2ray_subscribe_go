package worker

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
)

type Init struct {
}

func (p *Init) Init() (needRun bool, err error) {
	return true, InitWorker()
}

func (p *Init) WaitFinish() {

}

func (p *Init) Name() string {
	return "worker"
}

func InitWorker() error {
	finishC := make(chan bool, 1)

	err := InitTohru()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = InitProxy(finishC)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	select {
	case <-finishC:
		log.Debugf("init proxy success")
	case <-time.After(time.Second * 10):
		log.Warn("proxy start timeout")
	}

	if !conf.Config.Crawler.Enable && !conf.Config.ProxyCheck.Enable {
		return nil
	}

	beforeFunChan := make(chan func(), 10)
	if conf.Config.Crawler.Enable {
		log.Info("register crawler")

		beforeFunChan <- func() {
			err := crawler()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}

		// c.AddFunc(gron.Every(xtime.Minute*10), func() {
		// 	err := crawler()
		// 	if err != nil {
		// 		log.Errorf("err:%v", err)
		// 	}
		// })
	} else {
		log.Debugf("crawler not start")
	}

	if conf.Config.ProxyCheck.Enable {
		log.Info("register proxy check")

		beforeFunChan <- func() {
			err = checkProxyNode()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}

		// c.AddFunc(gron.Every(xtime.Minute*10), func() {
		// 	err := checkProxyNode(proxyCheck)
		// 	if err != nil {
		// 		log.Errorf("err:%v", err)
		// 	}
		// })

		// sched.Schedule().EverySingle().Friday().At("00:00").Do(func() {
		// 	err := db.Db.
		// 		Where("death_count > ?", 20).
		// 		Updates(map[string]interface{}{
		// 			"death_count": "death_count - 10",
		// 		}).Error
		// 	if err != nil {
		// 		log.Errorf("err:%v", err)
		// 		return
		// 	}
		// })

	} else {
		log.Debugf("proxy check not start")
	}

	beforeFunChan <- nil

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for {
			f := <-beforeFunChan
			if f == nil {
				break
			}

			f()
		}

		close(beforeFunChan)

		crawlerTimer := time.NewTimer(time.Minute)
		defer crawlerTimer.Stop()
		proxyCheckTimer := time.NewTimer(time.Minute)
		defer proxyCheckTimer.Stop()

		for {
			select {
			case <-crawlerTimer.C:
				if conf.Config.Crawler.Enable {
					err = crawler()
					if err != nil {
						log.Errorf("err:%v", err)
					}
					crawlerTimer.Reset(time.Minute)
				}
			case <-proxyCheckTimer.C:
				if conf.Config.ProxyCheck.Enable {
					err = crawler()
					if err != nil {
						log.Errorf("err:%v", err)
					}
					proxyCheckTimer.Reset(time.Minute)
				}
			case <-sigCh:
				log.Warn("worker stop")
				return
			}
		}
	}()

	return nil
}
