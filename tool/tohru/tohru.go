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
	ConfigPath  string `arg:"-c,--config" help:"config file path"`
	Action      string `arg:"-s" help:"action:register\nchangepwd\ncheck" default:"check"`
	UserName    string `arg:"-u,--username" help:"user name" default:""`
	Password    string `arg:"-p,--password" help:"password" default:""`
	OldPassword string `arg:"-o,--old-password" help:"old password" default:""`
	NewPassword string `arg:"-n,--new-password" help:"new password" default:""`
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

	err = tohru.Tohru.Init()
	if err != nil {
		color.Red.Printf("err:%v\n", err)
		return
	}

	switch cmdArgs.Action {
	case "register":
		if cmdArgs.UserName == "" {
			color.Red.Printf("missed user name\n")
			return
		}

		if cmdArgs.Password == "" {
			color.Red.Printf("missed password\n")
			return
		}

		err = tohru.Tohru.Registration(cmdArgs.UserName, cmdArgs.Password)
		if err != nil {
			color.Red.Printf("err:%v\n", err)
			return
		}

		color.Green.Printf("register success")
		color.Green.Printf("user name:%v\n", cmdArgs.UserName)
		color.Green.Printf("password:%v\n", cmdArgs.Password)

	case "changepwd":
		if cmdArgs.UserName == "" {
			color.Red.Printf("missed user name\n")
			return
		}

		if cmdArgs.OldPassword == "" {
			color.Red.Printf("missed old password\n")
			return
		}

		if cmdArgs.NewPassword == "" {
			color.Red.Printf("missed new password\n")
			return
		}

		err = tohru.Tohru.ChangedPassword(cmdArgs.UserName, cmdArgs.OldPassword, cmdArgs.Password)
		if err != nil {
			color.Red.Printf("err:%v\n", err)
			return
		}

		color.Green.Printf("changed password success")
		color.Green.Printf("user name:%v\n", cmdArgs.UserName)
		color.Green.Printf("new password:%v\n", cmdArgs.Password)

	case "check":
		fallthrough
	default:
		err = tohru.Tohru.CheckUsable()
		if err != nil {
			color.Red.Printf("err:%v\n", err)
			return
		}
		color.Green.Printf("check success")
	}
}
