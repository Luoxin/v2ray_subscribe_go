package main

import (
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/eddieivan01/nic"
	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v10"
)

var validate = validator.New()

var cmdArgs struct {
	Url string `arg:"-u,--url" help:"url" validate:"required"`
}

func main() {
	arg.MustParse(&cmdArgs)
	err := validate.Struct(cmdArgs)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	checkUrl := cmdArgs.Url
	if strings.HasPrefix(checkUrl, "https://") {
		checkUrl = strings.Replace(checkUrl, "https://", "http://", 1)
	}

	resp, err := nic.Get(checkUrl, &nic.H{
		Timeout:           60,
		SkipVerifyTLS:     true,
		AllowRedirect:     true,
		DisableKeepAlives: true,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	if resp.StatusCode != 200 {
		log.Println(color.Yellow.Text("cannot connect"))
		return
	}

	log.Println(color.Green.Text("check success"))
}
