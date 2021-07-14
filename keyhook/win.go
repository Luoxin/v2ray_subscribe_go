// +build windows

package keyhook

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Luoxin/Eutamias/node"
	"github.com/Luoxin/Eutamias/parser"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/atotto/clipboard"
	"github.com/bluele/gcache"
	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
	log "github.com/sirupsen/logrus"
)

func InitKeyHook() error {
	cache := gcache.New(128).LRU().Build()
	copyFunc := func(event hook.Event) {
		content, err := clipboard.ReadAll()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}

		_, err = cache.Get(utils.Md5(content))
		if err == nil {
			return
		} else if err == gcache.KeyNotFoundError {

		} else {
			log.Errorf("err:%v", err)
			return
		}

		_ = cache.SetWithExpire(utils.Md5(content), true, time.Hour)

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
	robotgo.EventHook(hook.KeyDown, []string{"ctrl", "v"}, copyFunc)

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
