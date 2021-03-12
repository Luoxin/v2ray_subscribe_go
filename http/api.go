package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"subscribe/node"
	"subscribe/pac"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xxjwxc/ginrpc/api"

	"subscribe/conf"
	"subscribe/domain"
	"subscribe/proxies"
	"subscribe/title"
	"subscribe/utils"
)

func Version(c *fiber.Ctx) error {
	return c.SendString(conf.Version)
}

type Subscribe struct {
}

type AddNodeReq struct {
	NodeUrl string `json:"node_url" validate:"required"`
}

type AddNodeRsp struct {
}

// @Router /version [post,get]
func (*Subscribe) Version(c *api.Context) {
	c.String(http.StatusOK, conf.Version)
}

// @Router /subscription [get]
func (*Subscribe) Subscription(c *api.Context) {
	nodes, err := node.GetUsableNodeList(50)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	titleGen := title.NewProxyTitle()

	var nodeList []string
	nodes.Each(func(node *domain.ProxyNode) {
		if node.NodeDetail == nil {
			return
		}

		var buf string
		switch domain.ProxyNodeType(node.ProxyNodeType) {
		case domain.ProxyNodeType_ProxyNodeTypeVmess:
			b, err := utils.Base64DecodeStripped(strings.TrimPrefix(node.NodeDetail.Buf, "vmess://"))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			m := map[string]interface{}{}
			err = json.Unmarshal([]byte(b), &m)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			m["ps"] = titleGen.Get()

			//ui := uuid.New()
			//m["id"] = ui.String()

			x, err := json.Marshal(m)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			buf = "vmess://" + base64.URLEncoding.EncodeToString(x)
		default:
			buf = node.NodeDetail.Buf
		}

		nodeList = append(nodeList, buf)
	})

	x := base64.StdEncoding.EncodeToString([]byte(strings.Join(nodeList, "\n")))

	c.String(http.StatusOK, x)
}

// @Router /sub/clash [get]
func (*Subscribe) SubClash(c *api.Context) {
	const clashCacheKey = "sub_clash"

	var force bool
	{
		val, _ := c.GetQuery("force")
		if val == "1" || strings.ToLower(val) == "true" {
			force = true
		}
	}

	if !force {
		value, err := conf.Cache.Get(clashCacheKey)
		if err != nil {
			if err != gcache.KeyNotFoundError {
				log.Errorf("err:%v", err)
				return
			}
		} else {
			c.String(http.StatusOK, value.(string))
			return
		}
	}

	nodes, err := node.GetUsableNodeList(50)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	p := proxies.NewProxies()
	nodes.Each(func(node *domain.ProxyNode) {
		if node.NodeDetail == nil {
			return
		}

		p.AppendWithUrl(node.NodeDetail.Buf)
	})

	val := p.ToClashConfig()
	if val != "" {
		err = conf.Cache.SetWithExpire(clashCacheKey, val, time.Minute*5)
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}

	c.String(http.StatusOK, val)
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

// @Router /pac [post,get]
func (*Subscribe) Pac(c *api.Context) {
	c.String(http.StatusOK, pac.Get())
}
