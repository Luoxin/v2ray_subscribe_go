package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Dreamacro/clash/constant"
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

	stop := func() {
		log.Info("process stop")
	}

	restart(true)

	pac.InitPac()

	const restartInterval = time.Minute * 30
	restartTimer := time.NewTicker(restartInterval)
	checkTimer := time.NewTicker(time.Minute * 5)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-restartTimer.C:
			log.Info("restart proxy")
			restart(false)

		case <-checkTimer.C:
			var aliveCount, proxyCount float64

			proxyList := tunnel.Proxies()

			for proxyName, proxy := range proxyList {
				switch proxy.Type() {
				case constant.Direct:
					goto NEXT
				case constant.Reject:
					goto NEXT

				case constant.Shadowsocks:

				case constant.ShadowsocksR:

				case constant.Snell:

				case constant.Socks5:

				case constant.Http:

				case constant.Vmess:

				case constant.Trojan:

				case constant.Relay:
					goto NEXT
				case constant.Selector:
					goto NEXT
				case constant.Fallback:
					goto NEXT
				case constant.URLTest:
					goto NEXT
				case constant.LoadBalance:
					goto NEXT

				default:
					goto NEXT
				}

				proxyCount++
				log.Infof("%v(%v):%v", proxyName, proxy.Alive(), proxy.LastDelay())
				if !proxy.Alive() {
					continue
				}

				// if proxy.LastDelay() > 500 {
				// 	continue
				// }

				aliveCount++

			NEXT:
			}

			healthiness := aliveCount / proxyCount

			log.Infof("uesd proxies healthiness is %.2f%%", healthiness*100)
			if healthiness < 0.1 {
				restart(false)
				restartTimer.Reset(restartInterval)
			}

		case <-sigCh:
			stop()
			os.Exit(0)
			return
		}
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
