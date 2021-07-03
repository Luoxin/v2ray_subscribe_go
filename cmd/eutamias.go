package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	eutamias "github.com/Luoxin/Eutamias"
	log2 "github.com/Luoxin/Eutamias/log"
	"github.com/Luoxin/Eutamias/notify"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/alexflint/go-arg"
	"github.com/inconshreveable/go-update"
	"github.com/kardianos/service"
	log "github.com/sirupsen/logrus"
)

// netsh interface ipv4 show excludedportrange protocol=tcp
// https://superuser.com/questions/1486417/unable-to-start-kestrel-getting-an-attempt-was-made-to-access-a-socket-in-a-way
// https://gist.github.com/steeve/6905542
// goreleaser --snapshot --skip-publish --rm-dist

var cmdArgs struct {
	ConfigPath string `arg:"-c,--config" help:"config file path"`
	Action     string `arg:"-s" help:"install,uninstall,start,run"`
}

var UpdateUrl string

func doUpdate() {
	defer func() {
		c <- true
	}()
	if UpdateUrl == "" {
		return
	}

	log.Infof("check update...")

	client := http.DefaultClient
	client.Timeout = time.Minute

	resp, err := client.Get(UpdateUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}

var c = make(chan bool, 1)

func main() {
	log2.InitLog()

	go doUpdate()

	waitUpdate := func() {
		select {
		case <-c:
		case <-time.After(time.Second * 10):
		}
	}

	arg.MustParse(&cmdArgs)

	serConfig := &service.Config{
		Name:             eutamias.ServiceName,
		DisplayName:      eutamias.ServiceName,
		Description:      "一个可以自我维护的网络代理工具",
		Executable:       os.Args[0],
		WorkingDirectory: utils.GetExecPath(),
	}

	p := &Program{}
	s, err := service.New(p, serConfig)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	switch cmdArgs.Action {
	case "install":
		err = s.Install()
		if err != nil {
			fmt.Println("install err", err)
		} else {
			fmt.Println("install success")
		}
	case "start":
		waitUpdate()
		err = s.Start()
		if err != nil {
			fmt.Println("Start err", err)
		} else {
			fmt.Println("Start success")
		}
	case "uninstall", "remove":
		err = s.Stop()
		if err != nil {
			fmt.Println("install err", err)
		}
		err = s.Uninstall()
		if err != nil {
			fmt.Println("Uninstall err", err)
		} else {
			fmt.Println("Uninstall success")
		}
	case "restart":
		waitUpdate()
		err = s.Restart()
		if err != nil {
			fmt.Println("Restart err", err)
		} else {
			fmt.Println("Restart success")
		}
	case "stop":
		err = s.Stop()
		if err != nil {
			fmt.Println("Stop err", err)
		} else {
			fmt.Println("Stop success")
		}
	case "run":
		fallthrough
	default:
		waitUpdate()
		err = s.Run() // 运行服务
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}
}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	notify.Msg(fmt.Sprintf("%v: service starting", eutamias.ServiceName))
	log.Info("service starting...")
	go p.run(s)
	return nil
}

func (p *Program) run(s service.Service) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	err := eutamias.Init(cmdArgs.ConfigPath)
	if err != nil {
		log.Errorf("err:%v", err)
		notify.Msg(fmt.Sprintf("%v: service start fail", eutamias.ServiceName))
		return
	}

	defer func() {
		notify.Msg(fmt.Sprintf("%v: service start stop", eutamias.ServiceName))
	}()

	notify.Msg(fmt.Sprintf("%v: service started", eutamias.ServiceName))

	<-sigCh
}

func (p *Program) Stop(s service.Service) error {
	log.Info("service stop")
	return nil
}
