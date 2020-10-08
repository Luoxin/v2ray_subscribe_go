package main

import (
	"encoding/json"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const version = "0.0.0.1"

type Config struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`

	Debug bool `json:"debug"`

	DbAddr           string `json:"db_addr"`
	CheckPort        int64  `json:"check_port"`
	ProxiesCrawler   string `json:"proxies_crawler"`
	V2RayServicePath string `json:"v2ray_service_path"`
	EnableCheckAlive bool   `json:"enable_check_alive"`
	EnableCrawl      bool   `json:"enable_crawl"`
}

func initConfig() error {
	logFormatter := &nested.Formatter{
		FieldsOrder: []string{
			log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
			log.FieldKeyFunc, log.FieldKeyMsg,
		},
		TimestampFormat: time.RFC3339,
		HideKeys:        true,
		NoFieldsSpace:   true,
		//ShowFullLevel:         true,
		NoUppercaseLevel:      true,
		TrimMessages:          true,
		CallerFirst:           true,
		CustomCallerFormatter: nil,
	}

	log.SetFormatter(logFormatter)
	log.SetReportCaller(true)

	gin.DisableConsoleColor()
	gin.EnableJsonDecoderDisallowUnknownFields()
	gin.EnableJsonDecoderUseNumber()

	viper.SetConfigFile("./config.yaml")
	viper.AddConfigPath(".")

	// 配置一些默认值
	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "8080")

	err := viper.ReadInConfig()
	if err != nil {
		log.Errorf("err:%v", err)

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
		} else {
			// 配置文件被找到，但产生了另外的错误
		}
		return err
	}

	log.Infof("read config %+v", viper.AllSettings())

	j, err := json.Marshal(viper.AllSettings())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = json.Unmarshal(j, &s.Config)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("load conf: %+v", s)

	if s.Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return nil
}
