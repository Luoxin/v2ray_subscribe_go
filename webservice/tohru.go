package webservice

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elliotchance/pie/pie"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/session"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/proxynode"
	"github.com/Luoxin/Eutamias/tohru"
	"github.com/Luoxin/Eutamias/version"
)

var validate = validator.New()
var store *session.Store

const (
	SessionKeyUid     = "Uid"
	SessionKeyUserKey = "UserKey"
	SessionKeyUserIp  = "UserIp"
)

const slt = ".\";v&vm6vOyrS)Ew@ByjN1Er|=<9B~=PniQM4C4Ca=2V@%ZadNP\\Vd:I^}\\["

func genPassword(key, pwd string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s_%s_%s", key, slt, pwd))))
}

func getReq(c *fiber.Ctx, req interface{}) error {
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

func Registration(c *fiber.Ctx) error {
	var req tohru.RegistrationReq
	var rsp tohru.RegistrationRsp

	err := getReq(c, &req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var tohruFeed domain.TohruFeed
	err = db.Db.Where("user_id = ?", req.TohruKey).First(&tohruFeed).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {

		} else {
			log.Errorf("err:%v", err)
			return err
		}
	} else {
		return errors.New("tohru is already exist")
	}

	tohruFeed = domain.TohruFeed{
		UserId:       req.TohruKey,
		UserPassword: genPassword(req.TohruKey, req.TohruPassword),
		LastIp:       c.IP(),
	}

	log.Info(genPassword(req.TohruKey, req.TohruPassword))

	err = db.Db.Create(&tohruFeed).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return c.JSON(rsp)
}

func ChangePassword(c *fiber.Ctx) error {
	var req tohru.ChangePasswordReq
	var rsp tohru.ChangePasswordRsp

	err := getReq(c, &req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var tohruFeed domain.TohruFeed
	err = db.Db.Where("user_id = ?", req.TohruKey).First(&tohruFeed).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("user not found")
			return c.SendStatus(403)
		}
		log.Errorf("err:%v", err)
		return err
	}

	if tohruFeed.UserPassword != genPassword(req.TohruKey, req.OldTohruPassword) {
		log.Error("password is fail")
		return c.SendStatus(403)
	}

	if req.NewTohruPassword != req.OldTohruPassword {
		err = db.Db.Model(&tohruFeed).Where("id = ?", tohruFeed.Id).
			Updates(map[string]interface{}{
				"user_password": genPassword(req.TohruKey, req.NewTohruPassword),
			}).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return c.JSON(rsp)
}

func CheckUsable(c *fiber.Ctx) error {
	var req tohru.CheckUsableReq
	var rsp tohru.CheckUsableRsp

	err := getReq(c, &req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	hello, err := conf.Ecc.ECCDecrypt(req.Hello)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var userInfo tohru.UserInfo
	err = json.Unmarshal([]byte(hello), &userInfo)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if userInfo.Hello != tohru.Hello {
		log.Error("hello not match")
		return c.SendStatus(403)
	}

	var tohruFeed domain.TohruFeed
	err = db.Db.Where("user_id = ?", userInfo.TohruKey).First(&tohruFeed).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("user not found")
			return c.SendStatus(403)
		}
		log.Errorf("err:%v", err)
		return err
	}

	log.Info(genPassword(userInfo.TohruKey, userInfo.TohruPassword))

	if tohruFeed.UserPassword != genPassword(userInfo.TohruKey, userInfo.TohruPassword) {
		log.Error("password is fail")
		return c.SendStatus(403)
	}

	if tohruFeed.LastIp != c.IP() {
		db.Db.Model(&domain.TohruFeed{}).
			Where("id = ?", tohruFeed.Id).Updates(map[string]interface{}{
			"last_ip": c.IP(),
		})
	}

	sess, err := store.Get(c)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = sess.Regenerate()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	rsp.Token = sess.ID()

	err = sess.Destroy()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	sess.Set(SessionKeyUid, tohruFeed.Id)
	sess.Set(SessionKeyUserKey, userInfo.TohruKey)
	sess.Set(SessionKeyUserIp, c.IP())

	// save session
	err = sess.Save()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	rsp.Version = version.Version
	return c.JSON(rsp)
}

func SyncNode(c *fiber.Ctx) error {
	var req tohru.SyncNodeReq
	var rsp tohru.SyncNodeRsp

	err := getReq(c, &req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	b, err := conf.Ecc.ECCDecrypt(req.NodeList)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var nodeUrlList pie.Strings
	err = json.Unmarshal([]byte(b), &nodeUrlList)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var contribution int
	nodeUrlList.Each(func(s string) {
		isNew, err := proxynode.AddNodeWithUrl(s)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		if isNew {
			contribution++
		}
	})

	if contribution > 0 {
		log.Info(contribution)
		sess, err := store.Get(c)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = db.Db.Model(&domain.TohruFeed{}).
			Where("id = ?", sess.Get(SessionKeyUid)).
			Updates(map[string]interface{}{
				"up_count": gorm.Expr("up_count + ?", contribution),
			}).Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	rsp.NodeList, err = proxynode.GetNode4Tohru(10 + contribution)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return c.JSON(rsp)
}

func registerTohru(app fiber.Router) error {
	app.Use("/", func(c *fiber.Ctx) error {
		var whiteList = pie.Strings{
			"/tohru/CheckUsable",
			"/tohru/Registration",
			"/tohru/ChangePassword",
		}

		path := c.Path()

		if whiteList.Any(func(value string) bool {
			return value == path
		}) {
			return c.Next()
		}

		sess, err := store.Get(c)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		if sess.Get(SessionKeyUserIp) != c.IP() {
			return c.SendStatus(403)
		}

		return c.Next()
	})

	app.Use("/Registration", basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "admin",
		},
		Authorizer: func(user, pass string) bool {
			if user == "admin" && pass == "admin" {
				return true
			}
			return false
		},
	}))

	app.Post("/ChangePassword", ChangePassword)
	app.Post("/Registration", Registration)

	app.Post("/CheckUsable", CheckUsable)
	app.Post("/SyncNode", SyncNode)

	return nil
}
