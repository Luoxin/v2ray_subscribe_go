package db

import (
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"subsrcibe/domain"
)

var Db *gorm.DB

func InitDb(dbAddr string) error {
	log.Infof("connect to database %v", dbAddr)

	addrList := strings.Split(dbAddr, "://")
	if len(addrList) < 2 {
		log.Errorf("Wrong database address")
		return errors.New("invalid args")
	}

	var d gorm.Dialector
	switch strings.ToLower(addrList[0]) {
	case "sqlite":
		d = sqlite.Open(strings.Join(addrList[1:], ""))
	case "mysql":
		d = mysql.Open(strings.Join(addrList[1:], ""))
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
		&domain.CrawlerConf{},
		&domain.ProxyNode{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	Db = db.Debug()

	return nil
}