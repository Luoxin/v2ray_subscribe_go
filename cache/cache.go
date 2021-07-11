package cache

// 不支持分布式

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/utils/json"
	"github.com/garyburd/redigo/redis"
	"github.com/shomali11/xredis"
	// "github.com/philippgille/gokv"
	// "github.com/philippgille/gokv/encoding"
	// "github.com/philippgille/gokv/gomap"
	// "github.com/philippgille/gokv/redis"
	log "github.com/sirupsen/logrus"
)

// var storage gokv.Store

var client *xredis.Client

func InitCache() (err error) {
	cacheConf := conf.Config.Cache

	switch cacheConf.Typ {
	case "redis":
		db, _ := strconv.ParseInt(cacheConf.DB, 10, 32)

		client = xredis.NewClient(&redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", cacheConf.Address,
					redis.DialDatabase(int(db)),
					redis.DialPassword(cacheConf.Password),
					redis.DialConnectTimeout(time.Second*3),
					redis.DialReadTimeout(time.Second*3),
					redis.DialWriteTimeout(time.Second*3),
					redis.DialKeepAlive(time.Minute),
				)
			},
		})
		// storage, err = redis.NewClient(redis.Options{
		// 	Address:  cacheConf.Address,
		// 	Password: cacheConf.Password,
		// 	DB:       int(db),
		// })
		// if err != nil {
		// 	log.Errorf("err:%v", err)
		// 	return err
		// }
	default:
		log.Warnf("not enable client")
		// storage = gomap.NewStore(gomap.Options{
		// 	Codec: encoding.JSON,
		// })
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-sigCh
		// _ = storage.Close()
		if client != nil {
			_ = client.Close()
		}
		log.Info("cache closed")
	}()

	return nil
}

var (
	ErrNotFound = errors.New("cache key not found")
)

func genKey(k string) string {
	return "eutamias_cache_" + k
}

func Get(k string, v interface{}) error {
	if client == nil {
		return ErrNotFound
	}

	value, _, err := client.Get(genKey(k))
	if err == redis.ErrNil {
		return ErrNotFound
	} else if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	v = value

	// found, err := storage.Get(genKey(k), v)
	// if err != nil {
	// 	log.Errorf("err:%v", err)
	// 	return err
	// }
	//
	// if !found {
	// 	return ErrNotFound
	// }
	return nil
}

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

func toJson(v interface{}) string {
	var value string
	switch x := v.(type) {
	case string:
		value = x
	case []byte:
		value = string(x)
	case int, int8, int16, int32, int64:
		value = fmt.Sprintf("%v", x)
	case uint, uint8, uint16, uint32, uint64:
		value = fmt.Sprintf("%v", x)
	case float32, float64:
		value = fmt.Sprintf("%v", x)
	default:
		b, _ := json.Marshal(v)
		value = string(b)
	}

	return value
}

func Set(k string, v interface{}) error {
	if client == nil {
		return nil
	}

	_, err := client.Set(genKey(k), toJson(v))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// return storage.Set(genKey(k), v)
	return nil
}

func SetEx(k string, v interface{}, timeout time.Duration) error {
	if client == nil {
		return nil
	}

	_, err := client.SetEx(genKey(k), toJson(v), int64(timeout.Seconds()))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// return storage.Set(genKey(k), v)
	return nil
}

func Del(k string) error {
	if client == nil {
		return nil
	}

	_, err := client.Del(genKey(k))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
	// return storage.Delete(genKey(k))
}

func Incr(k string) error {
	if client == nil {
		return nil
	}

	_, err := client.Incr(genKey(k))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func IncrEx(k string, timeout time.Duration) error {
	if client == nil {
		return nil
	}

	_, err := client.Incr(genKey(k))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = Expire(k, timeout)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func Expire(k string, timeout time.Duration) error {
	if client == nil {
		return nil
	}

	_, err := client.Expire(genKey(k), int64(timeout.Seconds()))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
