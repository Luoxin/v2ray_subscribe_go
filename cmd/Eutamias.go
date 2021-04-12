package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	arg "github.com/alexflint/go-arg"

	"github.com/Luoxin/Eutamias"
)

// netsh interface ipv4 show excludedportrange protocol=tcp
// https://superuser.com/questions/1486417/unable-to-start-kestrel-getting-an-attempt-was-made-to-access-a-socket-in-a-way
// https://gist.github.com/steeve/6905542
// goreleaser --snapshot --skip-publish --rm-dist

func main() {
	var cmdArgs struct {
		ConfigPath string `env:"C" arg:"-c,--config" help:"config file path"`
	}

	arg.MustParse(&cmdArgs)

	err := subscribe.Init(cmdArgs.ConfigPath)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

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