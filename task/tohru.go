package task

import (
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/tohru"
)

func InitTohru() error {
	err := tohru.Tohru.Init()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = tohru.Tohru.Start()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
