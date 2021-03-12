package http

import (
	"github.com/gofiber/fiber/v2"

	"subscribe/node"
	"subscribe/pac"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"subscribe/conf"
	"subscribe/domain"
	"subscribe/proxies"
)

func Version(c *fiber.Ctx) error {
	return c.SendString(conf.Version)
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

type Subscribe struct {
}

type AddNodeReq struct {
	NodeUrl string `json:"node_url" validate:"required"`
}

type AddNodeRsp struct {
}

// @Router /addnode [post,get]
func (*Subscribe) AddNode(c *gin.Context, req *AddNodeReq) (*AddNodeRsp, error) {
	var rsp AddNodeRsp

	err := node.AddNode(req.NodeUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
