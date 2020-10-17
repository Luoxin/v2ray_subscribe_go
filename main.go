package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/xxjwxc/ginrpc"
)

func main() {
	err := initState()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	//// swagger
	//myswagger.SetHost("https://localhost:8080")
	//myswagger.SetBasePath("gmsec")
	//myswagger.SetSchemes(true, false)
	//// -----end --

	base := ginrpc.New(ginrpc.WithDebug(s.Config.Debug), ginrpc.WithOutDoc(true))

	router := gin.Default()

	base.Register(router, new(SubscribeApi)) // 对象注册 like(go-micro)
	err = router.Run(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))

	//r := gin.Default()
	//err = registerRouting(r)
	//if err != nil {
	//	log.Errorf("err:%v", err)
	//	return
	//}

	//err = r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
