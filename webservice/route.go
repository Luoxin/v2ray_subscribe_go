package webservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/v2ray_subscribe_go/utils"
)

func registerRouting4Sub(sub fiber.Router) error {
	sub.Use("/",
		cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				refresh := c.Query("refresh")
				return refresh == "1" || strings.ToLower(refresh) == "true"
			},
			Expiration:   time.Minute * 5,
			CacheControl: true,
			// Storage:      storage,
		}))
	sub.Get("/v2ray/", timeout.New(SubV2ray, time.Minute))
	sub.Post("/v2ray/", timeout.New(SubV2ray, time.Minute))
	sub.Get("/clash/", timeout.New(SubClash, time.Minute))
	sub.Post("/clash/", timeout.New(SubClash, time.Minute))
	return nil
}

func registerRouting4Node(node fiber.Router) error {
	node.Post("/add", AddNode)
	node.Post("/addCrawlNode", AddCrawlerNode)
	node.Post("/list", NodeList)
	return nil
}

func registerRouting(app *fiber.App) error {
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Now().Sub(start)
		c.Response().Header.Set("X-Response-Time", fmt.Sprintf("%v", duration))
		return err
	})

	app.Group("/", func(c *fiber.Ctx) error {
		reqId := c.Context().ID()
		msg := fmt.Sprintf("<%v>[%s]%s %s %s", reqId, c.IP(), c.Method(), c.Path(), utils.ShortStr4Web(string(c.Body()), 400))
		log.Info(msg)

		err := c.Next()

		msg = fmt.Sprintf("<%v>%d %s", reqId, c.Response().StatusCode(),
			utils.ShortStr4Web(string(c.Response().Body()), 400))
		if err != nil {
			msg += fmt.Sprintf(" err:%v", err)
			log.Error(msg)
		} else {
			log.Info(msg)
		}

		return err
	})

	api := app.Group("/api/subscribe", func(c *fiber.Ctx) error {
		return c.Next()
	})

	api.Get("/version/", timeout.New(Version, time.Second))
	api.Post("/version/", timeout.New(Version, time.Second))
	api.Get("/pac/", timeout.New(Pac, time.Second))

	err := registerRouting4Sub(api.Group("/sub", func(c *fiber.Ctx) error {
		return c.Next()
	}))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = registerRouting4Node(api.Group("/node", func(c *fiber.Ctx) error {
		return c.Next()
	}))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = registerTohru(api.Group("/tohru"))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
