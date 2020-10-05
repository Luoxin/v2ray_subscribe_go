package main

import "github.com/gin-gonic/gin"

func registerRouting(r *gin.Engine) error {
	r.GET("/version", func(c *gin.Context) {
		c.String(200, "0.0.0.1")
	})

	return nil
}
