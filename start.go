package eutamias

import (
	"github.com/Luoxin/Eutamias/cache"
	"github.com/Luoxin/Eutamias/dns"
	"github.com/Luoxin/Eutamias/geolite"
	"github.com/Luoxin/Eutamias/initialize"
	"github.com/Luoxin/Eutamias/keyhook"
	"github.com/Luoxin/Eutamias/proxies"
	"github.com/pterm/pterm"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/webservice"
	"github.com/Luoxin/Eutamias/worker"
)

var ConfigFilePatch *string

var initList = []initialize.Initialize{
	&conf.Init{
		ConfigFilePatch: ConfigFilePatch,
	},
	&cache.Init{},
	&proxies.Init{},
	&db.Init{},
	&dns.InitClient{},
	&dns.InitServer{},
	&geolite.Init{},
	&worker.Init{},
	&keyhook.Init{},
	&webservice.Init{},
}

func Init(configFilePatch string) error {
	ConfigFilePatch = &configFilePatch

	p, err := pterm.DefaultProgressbar.
		WithTitle("Starting eutamias").
		WithShowElapsedTime(true).
		WithShowTitle(true).
		WithShowCount(true).
		WithTotal(len(initList)).
		Start()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for _, item := range initList {
		p.Title = "init " + item.Name()
		p.TitleStyle.Sprint(p.Title)

		needRun, err := item.Init()
		if err != nil {
			pterm.Error.Printfln("%v start fail:%v", item.Name(), err)
			return err
		}

		if !needRun {
			pterm.Warning.Printfln("%v does not need to be started", item.Name())
		} else {
			pterm.Success.Printfln("%v is started", item.Name())
		}

		p.Increment()
		item.WaitFinish()
	}

	return nil
}

func Start() {
	c := make(chan bool)

	err := Init("")
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

	<-c
}
