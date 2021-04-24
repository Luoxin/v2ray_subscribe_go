package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	eutamias "github.com/Luoxin/Eutamias"
	"github.com/Luoxin/Eutamias/conf"
	"github.com/alexflint/go-arg"
	"github.com/kardianos/service"
	"github.com/martinlindhe/notify"
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

const (
	serviceName = "eutamias"
)

func init() {
	conf.InitLog()
}

func main() {
	arg.MustParse(&cmdArgs)

	serConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: "eutamias service",
		Executable:  os.Args[0],
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
	case "uninstall", "remove":
		err = s.Uninstall()
		if err != nil {
			fmt.Println("Uninstall err", err)
		} else {
			fmt.Println("Uninstall success")
		}
	case "start":
		err = s.Start()
		if err != nil {
			fmt.Println("Start err", err)
		} else {
			fmt.Println("Start success")
		}
	case "restart":
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
		err = s.Run() // 运行服务
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	}
}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	notify.Notify(serviceName, serviceName, fmt.Sprintf("%v: service starting", serviceName), "")
	log.Info("service start")
	go p.run(s)
	return nil
}

func (p *Program) run(s service.Service) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	err := eutamias.Init(cmdArgs.ConfigPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	notify.Notify(serviceName, serviceName, fmt.Sprintf("%v: service started", serviceName), "")

	<-sigCh
}

func (p *Program) Stop(s service.Service) error {
	log.Info("service stop")
	return nil
}
