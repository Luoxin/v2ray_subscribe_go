package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/websocket/v2"
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
	logger := log.New()
	logger.SetLevel(log.InfoLevel)
	logger.Formatter = conf.LogFormatter
	app.Server().Logger = logger

	err := registerRouting(app)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	storage := sqlite3.New(sqlite3.Config{
		Database:   "fiber.vdb",
		Reset:      false,
		GCInterval: time.Hour,
	})
	app.Use(
		cache.New(cache.Config{
			Next: func(c *fiber.Ctx) bool {
				refresh := c.Query("refresh")
				return refresh == "1" || strings.ToLower(refresh) == "true"
			},
			Expiration:   time.Minute * 5,
			CacheControl: true,
			Storage:      storage,
		}),
		csrf.New(csrf.Config{
			KeyLookup:      "header:x-csrf-token",
			CookieName:     "csrf_",
			CookieSameSite: "Strict",
			Expiration:     time.Hour,
			Storage:        storage,
			KeyGenerator:   utils.UUID,
		}),
		compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		}),
		// rewrite.New(rewrite.Config{
		// 	Rules: map[string]string{
		// 		"/api/subscribe\\.subscription": "/api/subscribe/sub/v2ray",
		// 		"/api/subscribe\\.sub_clash":    "/api/subscribe/sub/clash",
		// 	},
		// }),
		// logger.New(logger.Config{
		// 	Format:       "[${time}] ${status} - ${latency} ${ip} ${ua} ${method} ${path} ${header} ${body}\n",
		// 	TimeFormat:   "15:04:05",
		// 	TimeZone:     "Local",
		// 	TimeInterval: 500 * time.Millisecond,
		// 	Output:       os.Stdout,
		// }),
	)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	err = app.Listen(fmt.Sprintf("%s:%d", conf.Config.HttpService.Host, conf.Config.HttpService.Port))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
