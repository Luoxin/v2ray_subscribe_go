package main

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type State struct {
	Config *Config
	Db     *gorm.DB
}

var s *State

func initState() error {
	s = &State{
		Config: &Config{},
	}

	err := initConfig()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = initDb()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = initCrawler()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
