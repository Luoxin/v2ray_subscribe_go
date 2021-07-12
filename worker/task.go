package worker

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/whiteshtef/clockwork"

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

	if conf.Config.IsTohru() {
		err := InitTohru()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	err := InitProxy(finishC)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	sched := clockwork.NewScheduler()

	select {
	case <-finishC:
		log.Debugf("init proxy success")
	case <-time.After(time.Minute * 10):
		log.Warn("proxy start timeout")
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

		sched.Schedule().EverySingle().Minute().Do(func() {
			err := crawler()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		})

		// c.AddFunc(gron.Every(xtime.Minute*10), func() {
		// 	err := crawler()
		// 	if err != nil {
		// 		log.Errorf("err:%v", err)
		// 	}
		// })
	} else {
		log.Warnf("crawler not start")
	}

	if conf.Config.ProxyCheck.Enable {
		log.Info("register proxy check")

		beforeFunChan <- func() {
			err = checkProxyNode()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		}

		sched.Schedule().EverySingle().Minute().Do(func() {
			err = checkProxyNode()
			if err != nil {
				log.Errorf("err:%v", err)
			}
		})

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
		log.Warnf("proxy check not start")
	}

	beforeFunChan <- nil
	go func() {
		for {
			f := <-beforeFunChan
			if f == nil {
				break
			}

			f()
		}
		sched.Run()
	}()

	return nil
}
