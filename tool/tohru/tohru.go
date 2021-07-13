package main

import (
	"github.com/Luoxin/Eutamias/conf"
	log2 "github.com/Luoxin/Eutamias/log"
	"github.com/Luoxin/Eutamias/tohru"
	"github.com/alexflint/go-arg"
	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"
)

var cmdArgs struct {
	ConfigPath string `arg:"-c,--config" help:"config file path"`
	Action     string `arg:"-s" help:"action" default:"check"`
}

func main() {
	arg.MustParse(&cmdArgs)
	log2.InitLog()
	log.SetLevel(log.DebugLevel)

	err := conf.InitConfig(cmdArgs.ConfigPath)
	if err != nil {
		color.Red.Printf("err:%v\n", err)
		return
	}

	switch cmdArgs.Action {
	case "check":
		fallthrough
	default:
		err = tohru.Tohru.Init()
		if err != nil {
			color.Red.Printf("err:%v\n", err)
			return
		}

		err = tohru.Tohru.CheckUsable()
		if err != nil {
			color.Red.Printf("err:%v\n", err)
			return
		}
		color.Green.Printf("check success")
	}
}
