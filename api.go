package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bluele/gcache"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xxjwxc/ginrpc/api"

	"subsrcibe/domain"
	"subsrcibe/utils"
)

func registerRouting(r *gin.Engine) error {
	//r.GET("/version", func(c *gin.Context) {
	//	c.String(200, version)
	//})

	return nil
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
	c.String(http.StatusOK, version)
}

// @Router /subscription [get]
func (*Subscribe) Subscription(c *api.Context) {
	nodes, err := GetUsableNodeList()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	titleGen := utils.NewProxyTitle()

	var nodeList []string
	ProxyNodeList(nodes).Each(func(node *domain.ProxyNode) {
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
		value, err := s.Cache.Get(clashCacheKey)
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

	nodes, err := GetUsableNodeList()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	val := utils.NewCoverSubscribe().Nodes2Clash(nodes)
	if val != "" {
		err = s.Cache.SetWithExpire(clashCacheKey, val, time.Minute*5)
		if err != nil {
			log.Errorf("err:%v", err)
		}
	}

	c.String(http.StatusOK, val)
}

// @Router /addnode [post,get]
func (*Subscribe) AddNode(c *gin.Context, node *AddNodeReq) (*AddNodeRsp, error) {
	var rsp AddNodeRsp

	err := addNode(node.NodeUrl, 0, 0)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

// @Router /pac [post,get]
func (*Subscribe) Pac(c *api.Context) {
}

func GetUsableNodeList() ([]*domain.ProxyNode, error) {
	var nodes []*domain.ProxyNode
	err := s.Db.Where("is_close = ?", false).
		Where("proxy_speed >= 0 ").
		// Where("proxy_node_type = 1").
		Where("available_count > 0 ").
		Where("proxy_network_delay >= 0").
		//Where("death_count < ?", 10).
		// Order("proxy_node_type").
		Order("available_count DESC").
		Order("proxy_speed DESC").
		Order("proxy_network_delay").
		Order("death_count").
		Order("last_crawler_at DESC").
		// Limit(50).
		Find(&nodes).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return nodes, err
}
