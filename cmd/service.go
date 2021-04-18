package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Luoxin/Eutamias/utils"
	"github.com/martinlindhe/notify"
	log "github.com/sirupsen/logrus"
	"github.com/takama/daemon"
	"golang.org/x/sys/windows/svc"
)

type Service struct {
	daemon.Daemon
}

func run() {
	for {
		utils.FileWrite("D:/TMP.LOG", "notify")
		notify.Notify("app name", "notice", "some text", "")
		time.Sleep(time.Second * 10)
	}
}
func (*Service) Run() {
	run()
}

func main() {
	service, err := daemon.New("at", fmt.Sprintf("%v", time.Now().Minute()), daemon.SystemDaemon)
	if err != nil {
		log.Fatal("Error: ", err)
	}

	if len(os.Args) > 1 {
		status, err := service.Remove()
		if err != nil {
			log.Fatal(status, "\nError: ", err)
		}

		// status, err = service.Install()
		// if err != nil {
		// 	log.Fatal(status, "\nError: ", err)
		// }
	} else {
		svc.Run("at", &winservice{})
	}
}

type winservice struct {
}

func (w *winservice) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	s <- svc.Status{
		State: svc.Running,
	}

	run()
	return
}
