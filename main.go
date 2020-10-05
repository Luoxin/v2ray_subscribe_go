package main

import (
	"brick/log"
	"github.com/gin-gonic/gin"
)

func main() {
	err := loadConfig()
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

	err = r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
