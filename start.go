package subsrcibe

import (
	log "github.com/sirupsen/logrus"

	"subsrcibe/conf"
	"subsrcibe/db"
	"subsrcibe/http"
	"subsrcibe/task"
)

func Start() {
	c := make(chan bool)

	err := conf.InitConfig()
	if err != nil {
		log.Fatalf("init config err:%v", err)
		return
	}

	err = db.InitDb(conf.Config.DbAddr)
	if err != nil {
		log.Fatalf("init config err:%v", err)
		return
	}

	log.Info("init conf success")

	err = task.InitWorker()
	if err != nil {
		log.Fatalf("init work err:%v", err)
		return
	}

	log.Info("init worker success")

	err = http.InitHttpService()
	if err != nil {
		log.Fatalf("init work err:%v", err)
		return
	}

	log.Info("init http service success")

	<-c
}

