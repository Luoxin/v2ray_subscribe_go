package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	mv2ray "github.com/v2fly/vmessping/miniv2ray"
	"github.com/v2fly/vmessping/vmess"
	"io/ioutil"
	"strings"
	"time"
	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	applog "v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	commlog "v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/infra/conf"
)

func main() {
	u := "vmess://eyJhZGQiOiJnei41Mm5vZGUueHl6IiwiYWlkIjoiMiIsImhvc3QiOiJoay5qZC5jb20iLCJpZCI6ImVkZDExODEwLTJkN2YtMzFlOC1hYzhjLTNmMjRmNjk4NDk4ZiIsIm5ldCI6IndzIiwicGF0aCI6Ii92MnJheSIsInBvcnQiOjI4MDI1LCJwcyI6IlNTUlRPT0wuQ09NIiwidGxzIjoiIiwidHlwZSI6Im5vbmUiLCJ2IjoiMiJ9\n"

	check(u,
		"")
}

func check(buf, checkUrl string) error {
	buf = strings.TrimSuffix(buf, "\n")

	if checkUrl == "" {
		checkUrl = "http://cachefly.cachefly.net/10mb"
	}

	var speed float64
	var delay time.Duration

	server, err := StartV2Ray(buf, false)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if err = server.Start(); err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer server.Close()

	client, err := mv2ray.CoreHTTPClient(server, 60*time.Second)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	defer client.CloseIdleConnections()

	before := time.Now()
	resp, err := client.Get(checkUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	delay = time.Now().Sub(before)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("err:%v", err)
	}
	defer resp.Body.Close()

	if len(body) > 0 {
		speed = float64(len(body)) / delay.Seconds()
	}

	log.Infof("%v(%v) speed %v, delay %v", buf, checkUrl, speed, delay)
	return nil
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
