// //go:generate goversioninfo -icon=resource/icon.ico -manifest=resource/goversioninfo.exe.manifest -64
// https://www.freebuf.com/sectool/246285.html
// https://docs.microsoft.com/zh-cn/dotnet/framework/tools/signtool-exe
// https://blog.csdn.net/wangshubo1989/article/details/50849914
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xxjwxc/ginrpc"
	"github.com/xxjwxc/ginrpc/api"
)

func main() {
	err := initState()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	log.Info("init success")

	//// swagger
	//myswagger.SetHost("https://localhost:8080")
	//myswagger.SetBasePath("gmsec")
	//myswagger.SetSchemes(true, false)
	//// -----end --

	base := ginrpc.New(
		ginrpc.WithCtx(api.NewAPIFunc),
		ginrpc.WithDebug(s.Config.Debug),
		ginrpc.WithOutDoc(true),
	)

	router := gin.Default()
	group := router.Group("/api")

	//group.GET("/version", base.HandlerFunc(Version))
	//group.POST("/version", base.HandlerFunc(Version))
	//
	//group.GET("/subscription", base.HandlerFunc(Subscription))
	//group.POST("/subscription", base.HandlerFunc(Subscription))

	base.Register(group, new(Subscribe))

	err = router.Run(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
