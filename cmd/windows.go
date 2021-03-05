package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Dreamacro/clash/hub"
	"github.com/Dreamacro/clash/hub/executor"
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

	for {
		select {
		case <-time.After(time.Minute * 30):
			restart(false)
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
