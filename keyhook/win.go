// +build windows

package keyhook

import (
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Luoxin/Eutamias/node"
	"github.com/Luoxin/Eutamias/parser"
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"
)

func InitKeyHook() error {
	copyFunc := func(event hook.Event) {
		content, err := clipboard.ReadAll()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		log.Infof("get %v from clipboard", content)

		parser.NewFuzzyMatchingParser().ParserText(content).Filter(func(s string) bool {
			return strings.Contains(s, "://")
		}).Each(func(nodeUrl string) {
			_, err = node.AddNodeWithUrl(nodeUrl)
			if err != nil {
				log.Errorf("link:%s, err:%v", nodeUrl, err)
				return
			}
		})
	}

	robotgo.EventHook(hook.KeyDown, []string{"ctrl", "c", "c"}, copyFunc)
	robotgo.EventHook(hook.KeyDown, []string{"ctrl", "x"}, copyFunc)

	s := robotgo.EventStart()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for {
			select {
			case <-robotgo.EventProcess(s):
			case <-sigCh:
				return
			}
		}
	}()

	return nil
}
