package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerRouting(r *gin.Engine) error {
	//r.GET("/version", func(c *gin.Context) {
	//	c.String(200, version)
	//})

	return nil
}

type SubscribeApi struct {
}

// Hello 带注解路由(参考beego形式)
// @Router /version [post,get]
func (s *SubscribeApi) Version(c *gin.Context) {
	c.JSON(http.StatusOK, version)
}
