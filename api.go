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

type VersionReq struct {
}

type VersionRsp struct {
}

type SubscriptionReq struct {
}

type SubscriptionRsp struct {
}

// Hello Annotated route (bese on beego way)
// @Router /version [post,get]
func (*Subscribe) Version(c *api.Context) {
	c.String(http.StatusOK, version)
}

// Hello Annotated route (bese on beego way)
// @Router /subscription [post,get]
func (*Subscribe) Subscription(c *api.Context) {
	var nodes []*subscription.ProxyNode
	err := s.Db.Where("is_close = ?", false).
		Where("next_check_at < ?", utils.Now()).
		Where("proxy_speed > 0 ").
		Where("death_count < ?", 3).
		Order("proxy_speed DESC").
		Find(&nodes).Error
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

			x, err := json.Marshal(m)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			buf = string(x)
		default:
			buf = node.NodeDetail.Buf
		}

		nodeList = append(nodeList, buf)
	})

	x := base64.URLEncoding.EncodeToString([]byte(strings.Join(nodeList, "\n")))

	c.String(http.StatusOK, x)
}

func Version(c *api.Context, gin *VersionReq) (*VersionRsp, error) {
	c.String(http.StatusOK, version)
	return nil, nil
}

func Subscription(c *gin.Context, req *SubscriptionReq) (*SubscriptionRsp, error) {
	var nodes []*subscription.ProxyNode
	err := s.Db.Where("is_close = ?", false).
		Where("next_check_at < ?", utils.Now()).
		Where("death_count < ?", 3).
		Find(&nodes).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	var nodeList []string
	ProxyNodeList(nodes).Each(func(node *subscription.ProxyNode) {
		if node.NodeDetail == nil {
			return
		}

		nodeList = append(nodeList, node.NodeDetail.Buf)
	})

	x := base64.URLEncoding.EncodeToString([]byte(strings.Join(nodeList, "\n")))

	c.String(http.StatusOK, x)
	return nil, nil
}
