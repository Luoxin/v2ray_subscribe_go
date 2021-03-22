package webservice

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	"github.com/Luoxin/v2ray_subscribe_go/conf"
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
	app := fiber.New()

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
		// requestid.New(requestid.Config{
		// 	Header: "x-request-id",
		// 	Generator: func() string {
		// 		return strings.ReplaceAll(utils.UUIDv4(), "-", "")
		// 	},
		// 	ContextKey: "request-id",
		// }),
		// logger.New(logger.Config{
		// 	Next: nil,
		// 	// For more options, see the Config section
		// 	Format:   "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		// 	TimeZone: "Local",
		// 	Output:   os.Stdout,
		// }),
	)

	err = app.Listen(fmt.Sprintf("%s:%d", conf.Config.HttpService.Host, conf.Config.HttpService.Port))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
