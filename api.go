package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/xxjwxc/ginrpc/api"
	"net/http"
	"strings"
	"subsrcibe/subscription"
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

// @Router /subscription [post,get]
func (*Subscribe) Subscription(c *api.Context) {
	nodes, err := GetUsableNodeList()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	titleGen := NewProxyTitle()

	var nodeList []string
	ProxyNodeList(nodes).Each(func(node *subscription.ProxyNode) {
		if node.NodeDetail == nil {
			return
		}

		var buf string
		switch subscription.ProxyNodeType(node.ProxyNodeType) {
		case subscription.ProxyNodeType_ProxyNodeTypeVmess:
			b, err := base64Decode(strings.TrimPrefix(node.NodeDetail.Buf, "vmess://"))
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

func GetUsableNodeList() ([]*subscription.ProxyNode, error) {
	var nodes []*subscription.ProxyNode
	err := s.Db.Where("is_close = ?", false).
		Where("proxy_speed >= 0 ").
		Where("proxy_node_type = 1").
		Where("available_count >= 0 ").
		Where("proxy_network_delay >= 0 ").
		Where("death_count < ?", 10).
		Order("proxy_node_type").
		Order("available_count DESC").
		Order("proxy_speed DESC").
		Order("proxy_network_delay DESC").
		Order("death_count").
		Order("last_crawler_at DESC").
		Limit(50).
		Find(&nodes).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return nodes, err
}
