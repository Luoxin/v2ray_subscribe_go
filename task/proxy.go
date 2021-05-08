package task

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub"
	"github.com/Dreamacro/clash/hub/executor"
	"github.com/Dreamacro/clash/tunnel"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/pac"
	"github.com/Luoxin/Eutamias/proxies"
	"github.com/Luoxin/Eutamias/utils"
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

	pwd, _ := os.Getwd()
	execPath := utils.GetExecPath()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		var isFirst = true

		homeDir, _ := os.UserHomeDir()
		clashConfigDir := filepath.Join(homeDir, ".config", "clash")

		clashGeoLiteFile := filepath.Join(clashConfigDir, "Country.mmdb")

		if !utils.FileExists(clashGeoLiteFile) {
			err := os.MkdirAll(clashConfigDir, 0777)
			if err != nil {
				log.Fatalf("err:%v", err)
				return
			}

			err = utils.CopyFile(filepath.Join(utils.GetConfigDir(), "GeoLite2.mmdb"), clashGeoLiteFile)
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
			var quantity = 50
			if isFirst {
				quantity = 10
			}

			config, count := proxies.GenClashConfig(quantity, !isFirst, isFirst)
			if count == 0 {
				log.Warnf("not found usabel proxies, skip update")
				return
			}

			log.Infof("get proxies %v", count)

			clashConf, err := executor.ParseWithBytes([]byte(config))
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			var options []hub.Option
			options = append(options, hub.WithExternalController(clashConf.General.ExternalController))

			{

				if utils.IsDir(filepath.Join(pwd, "./ui")) {
					log.Infof(filepath.Join(pwd, "./ui"))
					options = append(options, hub.WithExternalUI(filepath.Join(pwd, "./ui")))
				} else if utils.IsDir(filepath.Join(execPath, "./ui")) {
					log.Infof(filepath.Join(execPath, "./ui"))
					options = append(options, hub.WithExternalUI(filepath.Join(execPath, "./ui")))
				}
			}

			err = hub.Parse(options...)
			if err != nil {
				log.Errorf("err:%v", err) // TODO: fix:The system cannot find the file specified.
				return
			}

			executor.ApplyConfig(clashConf, force)
			isFirst = false
		}

		restart(true)
		finish()
		restart(true)

		pac.InitPac()

		rtf := time.NewTicker(time.Hour)
		rth := time.NewTicker(time.Minute * 45)
		rtm := time.NewTicker(time.Minute * 20)
		rtl := time.NewTicker(time.Minute * 5)

		getHealthiness := func() float64 {
			var proxyCount float64
			var less100, less500, less1000, less2000, alive float64
			//  5%       %20      %40       %20       %15

			proxyList := tunnel.Proxies()

			needSkip := func(proxy C.Proxy) bool {
				switch proxy.Type() {
				case constant.Direct, constant.Reject:
					return true

				case constant.Shadowsocks, constant.ShadowsocksR, constant.Snell,
					constant.Socks5, constant.Http, constant.Vmess, constant.Trojan:
					return false

				case constant.Relay, constant.Selector, constant.Fallback, constant.URLTest, constant.LoadBalance:
					return true

				default:
					return true
				}
			}

			for proxyName, proxy := range proxyList {
				if needSkip(proxy) {
					continue
				}

				proxyCount++
				log.Debugf("%v(%v):%v", proxyName, proxy.Alive(), proxy.LastDelay())
				if !proxy.Alive() {
					continue
				}

				delay := proxy.LastDelay()
				if delay < 100 {
					less100++
				}

				if delay < 500 {
					less500++
				}

				if delay < 1000 {
					less1000++
				}

				if delay < 2000 {
					less2000++
				}

				alive++
			}

			healthiness := (less100/proxyCount)*0.05 +
				(less500/proxyCount)*0.2 +
				(less1000/proxyCount)*0.4 +
				(less2000/proxyCount)*0.2 +
				(alive/proxyCount)*0.15

			log.Infof("uesd proxies healthiness is %.2f%%(%0.f|%0.f|%0.f|%0.f|%0.f|%0.f)",
				healthiness*100, less100, less500, less1000, less2000, alive, proxyCount)

			return healthiness
		}

		for {
			select {
			case <-rtf.C:
				if getHealthiness() < 0.9 {
					log.Info("restart proxy health lower then 0.9")
					restart(true)
				}

			case <-rth.C:
				if getHealthiness() < 0.7 {
					log.Info("restart proxy health lower then 0.7")
					restart(true)
				}

			case <-rtm.C:
				if getHealthiness() < 0.5 {
					log.Info("restart proxy health lower then 0.5")
					restart(false)
				}

			case <-rtl.C:
				if getHealthiness() < 0.2 {
					log.Info("restart proxy health lower then 0.2")
					restart(false)
				}
			case <-sigCh:
				log.Info("proxy stop")
				return
			}
		}
	}()

	return nil
}
