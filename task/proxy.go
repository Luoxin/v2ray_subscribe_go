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

	"github.com/luoxin/v2ray_subscribe_go/conf"
	"github.com/luoxin/v2ray_subscribe_go/domain"
	"github.com/luoxin/v2ray_subscribe_go/node"
	"github.com/luoxin/v2ray_subscribe_go/pac"
	"github.com/luoxin/v2ray_subscribe_go/proxies"
	"github.com/luoxin/v2ray_subscribe_go/utils"
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

		const restartInterval = time.Minute * 30
		restartTimer := time.NewTicker(restartInterval)
		checkTimer := time.NewTicker(time.Minute * 5)

		for {
			select {
			case <-restartTimer.C:
				log.Info("restart proxy")
				restart(true)

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

				log.Infof("uesd proxies healthiness is %.2f%%(%0.f/%0.f)", healthiness*100, aliveCount, proxyCount)
				if healthiness < 0.5 {
					restart(true)
					restartTimer.Reset(restartInterval)
				}

			case <-sigCh:
				stop()
				os.Exit(0)
				return
			}
		}
	}()

	return nil
}
