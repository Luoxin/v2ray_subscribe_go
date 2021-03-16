package webservice

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v10"

	"subscribe/conf"
	"subscribe/tohru"
	"subscribe/version"
)

var validate = validator.New()

func GetReq(c *fiber.Ctx, req interface{}) error {
	err := c.BodyParser(req)
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

type CheckUsableReq struct {
	// 当前Tohru的版本
	Version string `yaml:"version" json:"version" validate:"required"`
	Hello   string `yaml:"hello" json:"hello" validate:"required"`
}

type CheckUsableRsp struct {
	// 当前Kobayashi-san的版本
	Version string `yaml:"version" json:"version" validate:"required"`
}

func CheckUsable(c *fiber.Ctx) error {
	var req CheckUsableReq
	var rsp CheckUsableRsp

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

	rsp.Version = version.Version
	return c.JSON(rsp)
}

func registerTohru(app fiber.Router) error {
	app.Post("CheckUsable", CheckUsable)

	return nil
}
