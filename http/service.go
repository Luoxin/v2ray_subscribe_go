package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xxjwxc/ginrpc"
	"github.com/xxjwxc/ginrpc/api"

	"subscribe/conf"
)

func InitHttpService() error {
	if !conf.Config.HttpService.Enable {
		log.Warnf("http service not start")
		return nil
	}

	//// swagger
	//myswagger.SetHost("https://localhost:8080")
	//myswagger.SetBasePath("gmsec")
	//myswagger.SetSchemes(true, false)
	//// -----end --

	base := ginrpc.New(
		ginrpc.WithCtx(api.NewAPIFunc),
		ginrpc.WithDebug(conf.Config.Debug),
		ginrpc.WithOutDoc(true),
	)

	router := gin.Default()
	group := router.Group("/api")
	base.Register(group, new(Subscribe))

	go func() {
		err := router.Run(fmt.Sprintf("%s:%d", conf.Config.HttpService.Host, conf.Config.HttpService.Port))
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}()

	return nil
}
