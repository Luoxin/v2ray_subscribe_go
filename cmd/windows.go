package main

import (
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/tunnel"
	log "github.com/sirupsen/logrus"

	"subscribe"
	"subscribe/domain"
	"subscribe/http"
	"subscribe/proxies"
)

const ()

func main() {
	err := subscribe.Init()
	if err != nil {
		log.Errorf("err:%v", err)
	}

	nodes, err := http.GetUsableNodeList()
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

	clashConf, err := executor.ParseWithBytes([]byte(p.ToClashConfig()))
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	clashConf.General.Inbound.MixedPort = 7891

	executor.ApplyConfig(clashConf, true)

	proxyList := tunnel.Proxies()
	for _, proxy := range proxyList {
		log.Infof("%v %v %v", proxy.Name(), proxy.LastDelay(), proxy.Alive())
	}

	// a := app.New()
	// w := a.NewWindow("Hello")
	//
	// hello := widget.NewLabel("Hello Fyne!")
	// w.SetContent(container.NewVBox(
	// 	hello,
	// 	widget.NewButton("Hi!", func() {
	// 		hello.SetText("Welcome :)")
	// 	}),
	// ))

	// w.ShowAndRun()
}
