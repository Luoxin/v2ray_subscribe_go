package webservice

import (
	"github.com/gofiber/fiber/v2"

	"subscribe/node"
	"subscribe/pac"
	"subscribe/version"

	log "github.com/sirupsen/logrus"

	"subscribe/domain"
	"subscribe/proxies"
)

func Version(c *fiber.Ctx) error {
	return c.SendString(version.Version)
}

func SubV2ray(c *fiber.Ctx) error {
	nodes, err := node.GetUsableNodeList(50)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p := proxies.NewProxies()
	nodes.Each(func(node *domain.ProxyNode) {
		if node.NodeDetail == nil {
			return
		}

		p.AppendWithUrl(node.NodeDetail.Buf)
	})

	return c.SendString(p.ToV2rayConfig())
}

func SubClash(c *fiber.Ctx) error {
	nodes, err := node.GetUsableNodeList(50)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p := proxies.NewProxies()
	nodes.Each(func(node *domain.ProxyNode) {
		if node.NodeDetail == nil {
			return
		}

		p.AppendWithUrl(node.NodeDetail.Buf)
	})

	return c.SendString(p.ToClashConfig())
}

func Pac(c *fiber.Ctx) error {
	return c.SendString(pac.Get())
}

type AddNodeReq struct {
	NodeUrl string `json:"node_url" validate:"required"`
}

func AddNode(c *fiber.Ctx) error {
	var req AddNodeReq
	err := c.BodyParser(&req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if req.NodeUrl != "" {
		err := node.AddNode(req.NodeUrl)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return c.SendStatus(200)
}
