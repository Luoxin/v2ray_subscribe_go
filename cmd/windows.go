package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Dreamacro/clash/hub"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/tunnel"
	log "github.com/sirupsen/logrus"

	"subscribe"
	"subscribe/domain"
	"subscribe/http"
	"subscribe/pac"
	"subscribe/proxies"
	"subscribe/utils"
)

func main() {

	err := subscribe.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	restart := func(force bool) {
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

		var options []hub.Option
		options = append(options, hub.WithExternalController(clashConf.General.ExternalController))

		{
			pwd, _ := os.Executable()
			uiPath := filepath.Join(filepath.Dir(pwd), "./ui")
			if utils.IsDir(uiPath) {
				options = append(options, hub.WithExternalUI(uiPath))
			}
		}

		err = hub.Parse(options...)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		executor.ApplyConfig(clashConf, force)
	}

	restart(true)

	pac.InitPac()

	const restartInterval = time.Minute * 30
	restartTimer := time.NewTicker(restartInterval)
	checkTimer := time.NewTicker(time.Minute * 5)

	for {
		select {
		case <-restartTimer.C:
			log.Info("restart proxy")
			restart(false)
		case <-checkTimer.C:
			var aliveCount int
			for _, proxy := range tunnel.Proxies() {
				if !proxy.Alive() {
					continue
				}

				if proxy.LastDelay() > 500 {
					continue
				}

				aliveCount++
				if aliveCount > 5 {
					goto WAIT
				}
			}

			log.Info("restart proxy because proxy death")
			restart(true)
			restartTimer.Reset(restartInterval)
		}

	WAIT:
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
