package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Dreamacro/clash/adapters/outbound"
	log "github.com/sirupsen/logrus"
	"github.com/v2fly/vmessping/vmess"
	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	applog "v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	commlog "v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/infra/conf"

	"subsrcibe/subscription"
	"subsrcibe/utils"
)

func checkProxyNode() error {
	var nodes []*subscription.ProxyNode
	err := s.Db.Where("is_close = ?", false).
		Where("next_check_at < ?", utils.Now()).
		Where("death_count < ?", 50).
		Order("next_check_at").
		Find(&nodes).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if len(nodes) == 0 {
		log.Warnf("not found nodes need check")
		return nil
	}

	var w sync.WaitGroup
	c := make(chan *subscription.ProxyNode)
	stopC := make(chan bool)

	asyncCheck := func() {
		for {
			select {
			case node := <-c:
				run := func(node *subscription.ProxyNode) {
					defer w.Done()
					checkNode(node)
				}
				run(node)
			case <-stopC:
				return
			}
		}
	}

	const maxProcess = 50

	for i := 0; i < maxProcess; i++ {
		go asyncCheck()
	}

	startAt := time.Now()
	ProxyNodeList(nodes).
		Each(func(node *subscription.ProxyNode) {
			w.Add(1)
			c <- node
			log.Infof("add new node(%v) to check", node.Url)
		})

	w.Wait()

	for i := 0; i < maxProcess; i++ {
		stopC <- true
	}

	log.Infof("node check use %v", time.Now().Sub(startAt))

	return nil
}

func checkNode(node *subscription.ProxyNode) {
	log.Infof("wail check proxy for %+v", node.Url)
	defer log.Infof("check proxy finish,%v next exec at %v", node.Url, node.NextCheckAt)

	err := func() error {
		if node.NodeDetail == nil {
			node.IsClose = true
			return nil
		}

		proxyConfig := utils.ParseProxy(node.NodeDetail.Buf)
		proxyItem := make(map[string]interface{})
		j, err := json.Marshal(proxyConfig)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = json.Unmarshal(j, &proxyItem)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		proxyItem["name"] = "test proxy"

		proxy, err := outbound.ParseProxy(proxyItem)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second*60)

		// x, err := proxy.URLTest(ctx, "https://www.google.com/generate_204")
		delay, err := proxy.URLTest(ctx, "https://www.google.com")
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		node.ProxyNetworkDelay = float64(delay)
		node.AvailableCount++
		return nil
	}()
	if err != nil {
		log.Errorf("err:%v", err)
		node.DeathCount++
		node.AvailableCount = 0

		if node.DeathCount > 10 {
			node.ProxySpeed = -1
			node.ProxyNetworkDelay = -1
		}
	} else {
		node.DeathCount = 0
		log.Infof("check proxy %+v: speed:%v, delay:%v", node.Url, node.ProxySpeed, node.ProxyNetworkDelay)
	}

	node.NextCheckAt = node.CheckInterval + utils.Now()
	err = s.Db.Omit("node_detail", "url", "proxy_node_type").Save(node).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}

func StartV2Ray(vm string, debug bool) (*core.Instance, error) {

	loglevel := commlog.Severity_Error
	if debug {
		loglevel = commlog.Severity_Debug
	}

	lk, err := vmess.ParseVmess(vm)
	if err != nil {
		return nil, err
	}

	//fmt.Println("\n" + lk.DetailStr())
	ob, err := Vmess2Outbound(lk, true)
	if err != nil {
		return nil, err
	}
	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&applog.Config{
				ErrorLogType:  applog.LogType_Console,
				ErrorLogLevel: loglevel,
			}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	}

	commlog.RegisterHandler(commlog.NewLogger(commlog.CreateStderrLogWriter()))
	config.Outbound = []*core.OutboundHandlerConfig{ob}

	server, err := core.New(config)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func Vmess2Outbound(v *vmess.VmessLink, usemux bool) (*core.OutboundHandlerConfig, error) {

	out := &conf.OutboundDetourConfig{}
	out.Tag = "proxy"
	out.Protocol = "vmess"
	out.MuxSettings = &conf.MuxConfig{}
	if usemux {
		out.MuxSettings.Enabled = true
		out.MuxSettings.Concurrency = 8
	}

	p := conf.TransportProtocol(v.Net)
	s := &conf.StreamConfig{
		Network:  &p,
		Security: v.TLS,
	}

	switch v.Net {
	case "tcp":
		s.TCPSettings = &conf.TCPConfig{}
		if v.Type == "" || v.Type == "none" {
			s.TCPSettings.HeaderConfig = []byte(`{ "type": "none" }`)
		} else {
			pathb, _ := json.Marshal(strings.Split(v.Path, ","))
			hostb, _ := json.Marshal(strings.Split(v.Host, ","))
			s.TCPSettings.HeaderConfig = []byte(fmt.Sprintf(`
			{
				"type": "http",
				"request": {
					"path": %s,
					"headers": {
						"Host": %s
					}
				}
			}
			`, string(pathb), string(hostb)))
		}
	case "kcp":
		s.KCPSettings = &conf.KCPConfig{}
		s.KCPSettings.HeaderConfig = []byte(fmt.Sprintf(`{ "type": "%s" }`, v.Type))
	case "ws":
		s.WSSettings = &conf.WebSocketConfig{}
		s.WSSettings.Path = v.Path
		s.WSSettings.Headers = map[string]string{
			"Host": v.Host,
		}
	case "h2", "http":
		s.HTTPSettings = &conf.HTTPConfig{
			Path: v.Path,
		}
		if v.Host != "" {
			h := conf.StringList(strings.Split(v.Host, ","))
			s.HTTPSettings.Host = &h
		}
	}

	if v.TLS == "tls" {
		s.TLSSettings = &conf.TLSConfig{
			Insecure: true,
		}
		if v.Host != "" {
			s.TLSSettings.ServerName = v.Host
		}
	}

	out.StreamSetting = s
	oset := json.RawMessage(fmt.Sprintf(`{
  "vnext": [
    {
      "address": "%s",
      "port": %v,
      "users": [
        {
          "id": "%s",
          "alterId": %v,
          "security": "auto"
        }
      ]
    }
  ]
}`, v.Add, v.Port, v.ID, v.Aid))
	out.Settings = &oset
	return out.Build()
}
