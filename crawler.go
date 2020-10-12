package main

import (
	log "github.com/sirupsen/logrus"
	"time"
)

func initCrawler() error {
	if s.Config.DisableCrawl {
		log.Warnf("crawler disable")
		return nil
	}

	//go func() {
	log.Info("start crawler work...")
	//for {
	err := crawler()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	//}

	//}()

	return nil
}

func crawler() error {
	var crawlerList []*CrawlerConf
	err := s.Db.Where("is_closed = ? AND next_at < ?", true, time.Now().Unix()).Find(&crawlerList).Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info(crawlerList)

	return nil
}
