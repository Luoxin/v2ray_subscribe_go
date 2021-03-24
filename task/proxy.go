package task

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

	"github.com/Luoxin/v2ray_subscribe_go/conf"
	"github.com/Luoxin/v2ray_subscribe_go/domain"
	"github.com/Luoxin/v2ray_subscribe_go/node"
	"github.com/Luoxin/v2ray_subscribe_go/pac"
	"github.com/Luoxin/v2ray_subscribe_go/proxies"
	"github.com/Luoxin/v2ray_subscribe_go/utils"
)

func InitProxy(finishC chan bool) error {
	finish := func() {
		if finishC != nil {
			select {
			case finishC <- true:
			default:

			}
		}
	}

	if !conf.Config.Proxy.Enable {
		finish()
		return nil
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		var isFirst = true

		homeDir, _ := os.UserHomeDir()
		clashConfigDir := filepath.Join(homeDir, ".config", "clash")

		clashGeoLiteFile := filepath.Join(clashConfigDir, "Country.mmdb")

		if !utils.Exists(clashGeoLiteFile) {
			err := os.MkdirAll(clashConfigDir, 0777)
			if err != nil {
				log.Fatalf("err:%v", err)
				return
			}

			err = utils.CopyFile("./GeoLite2.mmdb", clashGeoLiteFile)
			if err != nil {
				log.Fatalf("err:%v", err)
				return
			}
		}

		// antPool, err := ants.NewPool(10)
		// if err != nil {
		// 	log.Errorf("err:%v", err)
		// 	return
		// }

		restart := func(force bool) {
			var quantity = -1
			if !isFirst {
				quantity = 10
			}

			nodes, err := node.GetUsableNodeList(quantity)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			p := proxies.NewProxies()

			// var w sync.WaitGroup
			nodes.Each(func(node *domain.ProxyNode) {
				// w.Add(1)
				// err = antPool.Submit(func() {
				// 	defer w.Done()
				p.AppendWithUrl(node.Url)
				// })
				// if err != nil {
				// 	log.Errorf("err:%v", err)
				// 	return
				// }
			})
			// w.Wait()

			if !isFirst {
				p = p.GetUsableList()
			}

			log.Infof("get proxies %v", p.Len())

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
			isFirst = false
		}

		stop := func() {
			log.Info("process stop")
		}

		restart(true)
		finish()
		restart(true)

		pac.InitPac()

		rtf := time.NewTicker(time.Hour)
		rtfh := time.NewTicker(time.Minute * 45)
		rtfm := time.NewTicker(time.Minute * 20)
		rtl := time.NewTicker(time.Minute * 5)

		getHealthiness := func() float64 {
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

			log.Infof("uesd proxies healthiness is %.2f%%(%0.f/%0.f)", healthiness*100, aliveCount, proxyCount)

			return healthiness
		}

		for {
			select {
			case <-rtf.C:
				log.Info("restart proxy forced")
				restart(true)

			case <-rtfh.C:
				if getHealthiness() < 0.9 {
					log.Info("restart proxy health lower then 0.9")
					restart(true)
				}

			case <-rtfm.C:
				if getHealthiness() < 0.6 {
					log.Info("restart proxy health lower then 0.6")
					restart(true)
				}

			case <-rtl.C:
				if getHealthiness() < 0.3 {
					log.Info("restart proxy health lower then 0.3")
					restart(true)
				}
			}
		}
	}()

	return nil
}
