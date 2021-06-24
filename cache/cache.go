package cache

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/gomap"
)

var Storage gokv.Store

func InitCache() error {
	cacheConf := conf.Config.Cache

	switch cacheConf.Typ {
	default:
		Storage = gomap.NewStore(gomap.Options{
			Codec: encoding.JSON,
		})
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-sigCh
		_ = Storage.Close()
	}()

	return nil
}
