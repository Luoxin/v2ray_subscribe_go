package webservice

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Luoxin/Eutamias/node"
	"github.com/Luoxin/Eutamias/pac"
	"github.com/Luoxin/Eutamias/version"

	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/proxies"
)

func Version(c *fiber.Ctx) error {
	return c.SendString(version.Version)
}

func SubV2ray(c *fiber.Ctx) error {
	nodes, err := node.GetUsableNodeList(50, true)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p := proxies.NewProxies()
	nodes.Each(func(node *domain.ProxyNode) {
		p.AppendWithUrl(node.Url)
	})

	return c.SendString(p.ToV2rayConfig())
}

func SubClash(c *fiber.Ctx) error {
	nodes, err := node.GetUsableNodeList(50, true)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p := proxies.NewProxies()
	nodes.Each(func(node *domain.ProxyNode) {
		p.AppendWithUrl(node.Url)
	})

	return c.SendString(p.ToClashConfig())
}

func Pac(c *fiber.Ctx) error {
	return c.SendString(pac.Get())
}

type AddNodeReq struct {
	NodeUrl string `json:"node_url" validate:"required"`
}

type AddCrawlerNodeReq struct {
	NodeUrl     string                   `json:"node_url" validate:"required"`
	CrawlerType domain.CrawlType         `json:"crawler_type"`
	Rule        *domain.CrawlerConf_Rule `json:"rule"`
}

func AddNode(c *fiber.Ctx) error {
	var req AddNodeReq
	err := c.BodyParser(&req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if req.NodeUrl != "" {
		_, err = node.AddNode(req.NodeUrl)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return c.SendStatus(200)
}

func AddCrawlerNode(c *fiber.Ctx) error {
	var req AddCrawlerNodeReq
	err := c.BodyParser(&req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if req.CrawlerType == domain.CrawlTypeNil {
		req.CrawlerType = domain.CrawlTypeFuzzyMatching
	}

	if req.NodeUrl != "" {
		err := node.AddCrawlerNode(req.NodeUrl, req.CrawlerType, req.Rule)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}

	return c.SendStatus(200)
}

func NodeList(ctx *fiber.Ctx) error {
	nodeList, err := node.GetUsableNodeList(100, true)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return ctx.JSON(nodeList)
}
