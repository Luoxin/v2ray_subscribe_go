package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	mv2ray "github.com/v2fly/vmessping/miniv2ray"
	"github.com/v2fly/vmessping/vmess"
	"io/ioutil"
	"strings"
	"subsrcibe/subscription"
	"subsrcibe/utils"
	"sync"
	"time"
	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	applog "v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	commlog "v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/infra/conf"
)

func initCheckProxy() error {
	if s.Config.DisableCheckAlive {
		log.Warn("check node is disable")
		return nil
	}

	go func() {
		log.Infof("check node starting...")

		for {
			err := checkProxyNode()
			if err != nil {
				log.Errorf("err:%v", err)
				continue
			}

			time.Sleep(time.Minute * 5)
		}

	}()

	return nil
}

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

	asyncCheck := func() {
		t := time.NewTicker(time.Second * 10)
		defer t.Stop()
		for {
			select {
			case node := <-c:
				checkNode(node)
				w.Done()
				t.Reset(time.Second * 10)
			case <-t.C:
				return
			}
		}
	}

	for i := 0; i < 10; i++ {
		go asyncCheck()
	}

	startAt := time.Now()
	ProxyNodeList(nodes).
		Each(func(node *subscription.ProxyNode) {
			w.Add(1)
			c <- node
		})

	w.Wait()
	log.Infof("node check use %v", time.Now().Sub(startAt))

	return nil
}

func checkNode(node *subscription.ProxyNode) {
	var speed, networkDelay float64
	err := func() error {
		if node.NodeDetail == nil {
			node.IsClose = true
			return nil
		}

		log.Infof("wail check proxy for %+v", node.Url)
		defer log.Infof("check proxy finish,%v next exec at %v", node.Url, node.NextCheckAt)

		switch subscription.ProxyNodeType(node.ProxyNodeType) {
		case subscription.ProxyNodeType_ProxyNodeTypeVmess:
			server, err := StartV2Ray(node.NodeDetail.Buf, false)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			if err = server.Start(); err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer server.Close()

			client, err := mv2ray.CoreHTTPClient(server, 10*time.Second)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			defer client.CloseIdleConnections()

			{
				// 存活性检测
				_, err := client.Get("https://www.google.com/generate_204")
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				} else {
					node.AvailableCount++
				}

			}

			{
				// 检测速度
				before := time.Now()
				resp, err := client.Get("http://cachefly.cachefly.net/10mb")
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				delay := time.Now().Sub(before)

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Errorf("err:%v", err)
					return err
				}
				defer resp.Body.Close()

				if len(body) > 0 {
					node.ProxySpeed = float64(len(body)) / delay.Seconds()
				}

				node.ProxyNetworkDelay = float64(delay.Milliseconds())
			}

		}

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
		node.ProxySpeed = speed
		node.ProxyNetworkDelay = networkDelay

		log.Infof("check proxy %+v: speed:%v, delay:%v", node.Url, speed, networkDelay)
	}

	node.NextCheckAt = node.CheckInterval + utils.Now()
	err = s.Db.Omit("node_detail", "death_count", "url", "proxy_node_type").Save(node).Error
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
