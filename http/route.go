package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	fiberUtils "github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"

	"subscribe/utils"
)

func registerRouting(app *fiber.App) error {
	base := app.Group("/", func(c *fiber.Ctx) error {
		reqId := strings.ReplaceAll(fiberUtils.UUIDv4(), "-", "")
		msg := fmt.Sprintf("<%s>[%s]%s %s %s", reqId, c.IP(), c.Method(), c.Path(), utils.ShortStr(string(c.Body()), 400))
		log.Info(msg)

		err := c.Next()

		msg = fmt.Sprintf("<%s>%d %s", reqId, c.Response().StatusCode(),
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

	api.All("version", timeout.New(Version, time.Second))
	api.All("pac", timeout.New(Pac, time.Second))

	sub := api.Group("sub/", func(c *fiber.Ctx) error {
		return c.Next()
	})
	sub.All("v2ray", timeout.New(SubV2ray, time.Minute))
	sub.All("clash", timeout.New(SubClash, time.Minute))

	return nil
}
