package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := initState()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	r := gin.Default()
	err = registerRouting(r)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = r.Run(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port)) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
