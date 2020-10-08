package main

import log "github.com/sirupsen/logrus"

type State struct {
	Config *Config
}

var s *State

func initState() error {
	s = &State{
		Config: &Config{},
	}

	err := initConfig()
	if err != nil {
		log.Error("err:%v", err)
		return err
	}

	return nil
}
