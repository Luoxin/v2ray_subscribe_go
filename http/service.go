package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	log "github.com/sirupsen/logrus"

	"subscribe/conf"
)

func InitHttpService() error {
	if !conf.Config.HttpService.Enable {
		log.Warnf("http service not start")
		return nil
	}

	// https://github.com/gofiber/fiber

	app := fiber.New()

	err := registerRouting(app)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			refresh := c.Query("refresh")
			return refresh == "1" || strings.ToLower(refresh) == "true"
		},
		Expiration:   time.Minute * 5,
		CacheControl: true,
	}))

	app.Use("/api", func(c *fiber.Ctx) error {
		fmt.Println("🥈 Second handler")
		return c.Next()
	})

	// GET /john
	app.Get("/:name", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Hello, %s 👋!", c.Params("name"))
		return c.SendString(msg) // => Hello john 👋!
	})

	// GET /john/75
	app.Get("/:name/:age", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("👴 %s is %s years old", c.Params("name"), c.Params("age"))
		return c.SendString(msg) // => 👴 john is 75 years old
	})

	// GET /dictionary.txt
	app.Get("/:file.:ext", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("📃 %s.%s", c.Params("file"), c.Params("ext"))
		return c.SendString(msg) // => 📃 dictionary.txt
	})

	// GET /flights/LAX-SFO
	app.Get("/flights/:from-:to", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("💸 From: %s, To: %s", c.Params("from"), c.Params("to"))
		return c.SendString(msg) // => 💸 From: LAX, To: SFO
	})

	// GET /api/register
	app.Get("/api/*", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("✋ %s", c.Params("*"))
		return c.SendString(msg) // => ✋ register
	})

	err = app.Listen(fmt.Sprintf("%s:%d", conf.Config.HttpService.Host, conf.Config.HttpService.Port))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
