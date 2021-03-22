package subscribe

import (
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/v2ray_subscribe_go/conf"
	"github.com/Luoxin/v2ray_subscribe_go/db"
	"github.com/Luoxin/v2ray_subscribe_go/task"
	"github.com/Luoxin/v2ray_subscribe_go/webservice"
)

func Init() error {
	err := conf.InitConfig()
	if err != nil {
		log.Fatalf("init config err:%v", err)
		return err
	}

	err = db.InitDb(conf.Config.Db.Addr)
	if err != nil {
		log.Fatalf("init db err:%v", err)
		return err
	}

	log.Info("init conf success")

	err = task.InitWorker()
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

	err := Init()
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

	<-c
}
