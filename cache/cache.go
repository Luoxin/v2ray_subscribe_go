package cache

// 不支持分布式

import (
	"errors"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/philippgille/gokv"
	"github.com/philippgille/gokv/encoding"
	"github.com/philippgille/gokv/gomap"
	"github.com/philippgille/gokv/redis"
	log "github.com/sirupsen/logrus"
)

var storage gokv.Store

func InitCache() (err error) {
	cacheConf := conf.Config.Cache

	switch cacheConf.Typ {
	case "redis":
		db, _ := strconv.ParseInt(cacheConf.DB, 10, 32)
		storage, err = redis.NewClient(redis.Options{
			Address:  cacheConf.Address,
			Password: cacheConf.Password,
			DB:       int(db),
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	default:
		storage = gomap.NewStore(gomap.Options{
			Codec: encoding.JSON,
		})
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-sigCh
		_ = storage.Close()
		log.Info("cache closed")
	}()

	return nil
}

func genKey(k string) string {
	return "eutamias_cache_" + k
}

func Get(k string, v interface{}) error {
	found, err := storage.Get(genKey(k), v)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if !found {
		return ErrNotFound
	}
	return nil
}

var ErrNotFound = errors.New("cache key not found")

func GetStr(k string) (string, error) {
	var v string
	return v, Get(k, &v)
}

func GetInt(k string) (int, error) {
	var v int
	return v, Get(k, &v)
}

func GetInt32(k string) (int32, error) {
	var v int32
	return v, Get(k, &v)
}

func GetInt64(k string) (int64, error) {
	var v int64
	return v, Get(k, &v)
}

func GetUint(k string) (uint, error) {
	var v uint
	return v, Get(k, &v)
}

func GetUint32(k string) (uint32, error) {
	var v uint32
	return v, Get(k, &v)
}

func GetUint64(k string) (uint64, error) {
	var v uint64
	return v, Get(k, &v)
}

func GetBool(k string) (bool, error) {
	var v bool
	return v, Get(k, &v)
}

func Set(k string, v interface{}) error {
	return storage.Set(genKey(k), v)
}

func Del(k string) error {
	return storage.Delete(genKey(k))
}
