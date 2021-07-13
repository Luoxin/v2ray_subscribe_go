package main

import (
	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/dns"
	log2 "github.com/Luoxin/Eutamias/log"
	"github.com/alexflint/go-arg"
	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"
)

var cmdArgs struct {
	ConfigPath string `arg:"-c,--config" help:"config file path"`
	Domain     string `arg:"-d,--domain" help:"domain"`
}

func main() {
	arg.MustParse(&cmdArgs)
	log2.InitLog()
	log.SetLevel(log.DebugLevel)
	log2.ShowConsole()

	err := conf.InitConfig(cmdArgs.ConfigPath)
	if err != nil {
		color.Red.Printf("err:%v\n", err)
		return
	}

	_, err = dns.InitDnsClient()
	if err != nil {
		color.Red.Printf("err:%v\n", err)
		return
	}

	ip := dns.LookupHostsFastestIp(cmdArgs.Domain)
	color.Green.Printf("%v %v\n", cmdArgs.Domain, ip)
}
