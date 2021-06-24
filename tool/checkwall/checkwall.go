package main

import (
	"net/url"

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

	u, err := url.Parse(cmdArgs.Url)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "http"
	case "":
		u.Scheme = "http"
	default:
	}

	resp, err := nic.Get(u.String(), &nic.H{
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
		color.Yellow.Println("cannot connect")
		return
	}

	color.Green.Println("check success")
}
