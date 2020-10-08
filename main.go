package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := initConfig()
	if err != nil {
		log.Error("err:%v", err)
		return
	}

	r := gin.Default()
	err = registerRouting(r)
	if err != nil {
		log.Error("err:%v", err)
		return
	}

	err = r.Run(fmt.Sprintf("%s:%d", s.Host, s.Port)) // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
