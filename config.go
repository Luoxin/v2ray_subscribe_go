package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	LogDebug         bool   `json:"log_debug"`
	DbUrl            string `json:"db_url"`
	CheckPort        int64  `json:"check_port"`
	ProxiesCrawler   string `json:"proxies_crawler"`
	V2RayServicePath string `json:"v2ray_service_path"`
	EnableCheckAlive bool   `json:"enable_check_alive"`
	EnableCrawl      bool   `json:"enable_crawl"`
}

func loadConfig() error {
	viper.SetConfigFile("./config.yaml")
	viper.AddConfigPath(".")

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

	return nil
}
