package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var s *Config

const version = "0.0.0.1"

type Config struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`

	Debug bool `json:"debug"`

	DbUrl            string `json:"db_url"`
	CheckPort        int64  `json:"check_port"`
	ProxiesCrawler   string `json:"proxies_crawler"`
	V2RayServicePath string `json:"v2ray_service_path"`
	EnableCheckAlive bool   `json:"enable_check_alive"`
	EnableCrawl      bool   `json:"enable_crawl"`
}

func initConfig() error {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		FullTimestamp:    false,
		TimestampFormat:  "2006-01-02 15:04:05.000",
		DisableSorting:   false,
		QuoteEmptyFields: false,
		DisableColors:    true,
		FieldMap: log.FieldMap{
			"@module": "v2ray_subscribe",
		},
	})

	log.SetReportCaller(true)
	gin.DisableConsoleColor()

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

	err = json.Unmarshal(j, &s)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("load conf: %+v", s)

	if s.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return nil
}
