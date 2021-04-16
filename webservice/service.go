package webservice

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/Luoxin/Eutamias/conf"
)

var storage fiber.Storage

func InitHttpService() error {
	if !conf.Config.HttpService.Enable {
		log.Warnf("http service not start")
		return nil
	}

	var err error
	storage, err = InitStorage(conf.Config.Db.Addr)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// https://github.com/gofiber/fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.SendStatus(500)
		},
		ServerHeader:  "",
		CaseSensitive: true,
		UnescapePath:  true,
		// ETag:                     true,
		ReadTimeout:              time.Minute * 5,
		WriteTimeout:             time.Minute * 5,
		CompressedFileSuffix:     ".gz",
		ProxyHeader:              "",
		DisableDefaultDate:       true,
		DisableHeaderNormalizing: true,
		ReduceMemoryUsage:        true,
	})

	err = InitWs(app)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	app.Server().ErrorHandler = func(ctx *fasthttp.RequestCtx, err error) {
		log.Errorf("%s err:%v", ctx.Request.String(), err)
	}

	err = registerRouting(app)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// storage := sqlite3.New(sqlite3.Config{
	// 	Database:   "fiber.vdb",
	// 	Reset:      false,
	// 	GCInterval: time.Hour,
	// })
	app.Use(
		// cache.New(cache.Config{
		// 	Next: func(c *fiber.Ctx) bool {
		// 		refresh := c.Query("refresh")
		// 		return refresh == "1" || strings.ToLower(refresh) == "true"
		// 	},
		// 	Expiration:   time.Minute * 5,
		// 	CacheControl: true,
		// 	// Storage:      storage,
		// }),
		csrf.New(csrf.Config{
			KeyLookup:      "header:x-csrf-token",
			CookieName:     "csrf_",
			CookieSameSite: "Strict",
			Expiration:     time.Hour,
			// Storage:        storage,
			KeyGenerator: utils.UUID,
		}),
		compress.New(compress.Config{
			Level: compress.LevelBestCompression,
		}),
		// logger.New(logger.Config{
		// 	Format:       "[${time}] ${status} - ${latency} ${ip} ${ua} ${method} ${path} ${header} ${body}\n",
		// 	TimeFormat:   "15:04:05",
		// 	TimeZone:     "Local",
		// 	TimeInterval: 500 * time.Millisecond,
		// 	Output:       os.Stdout,
		// }),
		requestid.New(requestid.Config{
			Header: "x-request-id",
			Generator: func() string {
				return strings.ReplaceAll(utils.UUIDv4(), "-", "")
			},
			ContextKey: "request-id",
		}),
		logger.New(logger.Config{
			Next: nil,
			// For more options, see the Config section
			Format:   "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
			TimeZone: "Local",
			Output:   os.Stdout,
		}),
		recover.New(recover.Config{
			Next: func(c *fiber.Ctx) bool {
				_ = c.SendStatus(500)
				return false
			},
			EnableStackTrace: false,
			StackTraceHandler: func(e interface{}) {
				log.Errorf("panic:%v", err)
			},
		}),
		// func(c *fiber.Ctx) {
		// 	_ = c.Status(500).JSON(map[string]interface{}{
		// 		"code": -1,
		// 		"msg":  "sys error",
		// 	})
		// 	_ = c.Status(404).JSON(map[string]interface{}{
		// 		"code": 404,
		// 		"msg":  "api not found",
		// 	})
		// },
		limiter.New(limiter.Config{
			Next: func(c *fiber.Ctx) bool {
				return c.IP() == "127.0.0.1"
			},
			Max: 40,
			KeyGenerator: func(c *fiber.Ctx) string {
				sess, err := store.Get(c)
				if err == nil {
					uid := sess.Get(SessionKeyUid)
					if uid != nil {
						return fmt.Sprintf("limiter_user_id_%v", uid)
					}
				}

				return fmt.Sprintf("limiter_user_ip_%s", c.IP())
			},
			Expiration: time.Minute,
		}),
	)

	err = app.Listen(fmt.Sprintf("%s:%d", conf.Config.HttpService.Host, conf.Config.HttpService.Port))
	if err != nil {
		log.Errorf("err:%v", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-sigCh
		err = app.Shutdown()
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}()

	return nil
}
