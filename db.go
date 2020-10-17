package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
	"subsrcibe/subscribe"
)

func initDb() error {
	log.Infof("connect to database %v", s.Config.DbAddr)

	addrList := strings.Split(s.Config.DbAddr, "://")
	if len(addrList) < 2 {
		log.Errorf("Wrong database address")
		return ErrInvalidArg
	}

	var d gorm.Dialector
	switch strings.ToLower(addrList[0]) {
	case "sqlite":
		d = sqlite.Open(strings.Join(addrList[1:], ""))
	default:
		return errors.New("unsupported database")
	}

	db, err := gorm.Open(d, &gorm.Config{
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		AllowGlobalUpdate:                        true,
		NamingStrategy: &schema.NamingStrategy{
			TablePrefix:   "subscribe_",
			SingularTable: true,
		},
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Infof("auto migrate tables")
	err = db.AutoMigrate(
		&subscribe.CrawlerConf{},
		&subscribe.ProxyNode{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	s.Db = db.Debug()

	return nil
}
