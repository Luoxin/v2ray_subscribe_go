package webservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	log "github.com/sirupsen/logrus"

	"subscribe/utils"
)

func registerRouting(app *fiber.App) error {
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Now().Sub(start)
		c.Response().Header.Set("X-Response-Time", fmt.Sprintf("%v", duration))
		return err
	})

	base := app.Group("/", func(c *fiber.Ctx) error {
		reqId := c.Context().ID()
		msg := fmt.Sprintf("<%v>[%s]%s %s %s", reqId, c.IP(), c.Method(), c.Path(), utils.ShortStr(string(c.Body()), 400))
		log.Info(msg)

		err := c.Next()

		msg = fmt.Sprintf("<%v>%d %s", reqId, c.Response().StatusCode(),
			strings.ReplaceAll(utils.ShortStr(string(c.Response().Body()), 400), "\n", "\\n"))
		if err != nil {
			msg += fmt.Sprintf(" err:%v", err)
			log.Error(msg)
		} else {
			log.Info(msg)
		}

		return err
	})

	api := base.Group("api/subscribe/", func(c *fiber.Ctx) error {
		return c.Next()
	})

	api.Get("version", timeout.New(Version, time.Second))
	api.Post("version", timeout.New(Version, time.Second))
	api.Get("pac", timeout.New(Pac, time.Second))

	sub := api.Group("sub/", func(c *fiber.Ctx) error {
		return c.Next()
	})
	sub.Get("v2ray", timeout.New(SubV2ray, time.Minute))
	sub.Post("v2ray", timeout.New(SubV2ray, time.Minute))
	sub.Get("clash", timeout.New(SubClash, time.Minute))
	sub.Post("clash", timeout.New(SubClash, time.Minute))

	node := api.Group("node/", func(c *fiber.Ctx) error {
		return c.Next()
	})
	node.Post("add", AddNode)

	return nil
}
