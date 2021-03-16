package webservice

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v10"

	"subscribe/conf"
	"subscribe/tohru"
	"subscribe/version"
)

var validate = validator.New()
var store = session.New(session.Config{
	Expiration:   time.Hour * 20,
	Storage:      nil,
	CookieName:   "x-tohru-id",
	KeyGenerator: nil,
})

func GetReq(c *fiber.Ctx, req interface{}) error {
	err := c.BodyParser(&req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = validate.Struct(req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func CheckUsable(c *fiber.Ctx) error {
	var req tohru.CheckUsableReq
	var rsp tohru.CheckUsableRsp

	err := GetReq(c, &req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	hello, err := conf.Ecc.ECCDecrypt(req.Hello)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if hello != tohru.Hello {
		c.Status(403)
		return errors.New("hello not match")
	}

	sess, err := store.Get(c)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// Destry session
	err = sess.Destroy()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// save session
	err = sess.Save()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	rsp.Version = version.Version
	rsp.Token = sess.ID()
	return c.JSON(rsp)
}

func SyncNode(c *fiber.Ctx) error {
	return c.Next()
}

func registerTohru(app fiber.Router) error {
	app.Use("/", func(c *fiber.Ctx) error {
		if c.Path() == "/api/subscribe/tohru/CheckUsable" {
			return c.Next()
		}

		return c.Next()
	})

	app.Post("/CheckUsable", CheckUsable)

	return nil
}
