package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const Version = "0.0.0.4"

type config struct {
	Host  string `json:"host"`
	Port  uint32 `json:"port"`
	Debug bool   `json:"debug"`

	DbAddr  string `json:"db_addr"`
	Proxies string `json:"proxies"`

	DisableCrawl       bool `json:"disable_crawl"`
	DisableCheckAlive  bool `json:"disable_check_alive"`
	DisableHttpService bool `json:"disable_http_service"`

	CrawlerInterval uint32 `json:"crawler_interval"`
	CheckInterval   uint32 `json:"check_interval"`
}

var Config config

func InitConfig() error {
	logFormatter := &nested.Formatter{
		FieldsOrder: []string{
			log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
			log.FieldKeyFunc, log.FieldKeyMsg,
		},
		CustomCallerFormatter: func(f *runtime.Frame) string {
			return fmt.Sprintf("(%s %s:%d)", f.Function, path.Base(f.File), f.Line)
		},
		TimestampFormat:  time.RFC3339,
		HideKeys:         true,
		NoFieldsSpace:    true,
		NoUppercaseLevel: true,
		TrimMessages:     true,
		CallerFirst:      true,
	}

	log.SetFormatter(logFormatter)
	log.SetReportCaller(true)

	// 可能存在的目录
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("./conf/")
	viper.AddConfigPath("../conf/")
	{
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Errorf("err:%v", err)
		} else {
			viper.AddConfigPath(filepath.Join(homeDir, "subscribe"))
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 配置一些默认值
	viper.SetDefault("port", 8080)
	viper.SetDefault("host", "127.0.0.1")

	viper.SetDefault("check_interval", 300)
	viper.SetDefault("crawler_interval", 3600)

	viper.SetDefault("db_addr", "sqlite://.subscribe.vdb?check_same_thread=false")

	viper.SetDefault("disable_crawl", false)
	viper.SetDefault("disable_check_alive", false)
	viper.SetDefault("disable_http_service", false)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("not found conf file, use default")
		} else if e, ok := err.(*os.PathError); ok {
			log.Fatalf("not find conf file in %s", e.Path)
		} else {
			log.Errorf("err:%v", err)
			return err
		}
	} else {
		// log.Infof("read config %+v", viper.AllSettings())

		// err = viper.Unmarshal(&Config)
		// if err != nil {
		// 	log.Errorf("err:%v", err)
		// 	return err
		// }

		j, err := json.Marshal(viper.AllSettings())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = json.Unmarshal(j, &Config)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		log.Infof("get config %+v", Config)
	}

	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return nil
}
