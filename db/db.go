package db

import (
	"errors"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/luoxin/v2ray_subscribe_go/subscribe/conf"
	"github.com/luoxin/v2ray_subscribe_go/subscribe/domain"
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
		d = mysql.New(mysql.Config{
			DSN: strings.Join(addrList[1:], ""),

			DefaultStringSize: 256,

			DontSupportRenameIndex: true,
		})
	default:
		return errors.New("unsupported database")
	}

	dbConfig := gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: &schema.NamingStrategy{
			TablePrefix:   "subscribe_",
			SingularTable: true,
		},
		FullSaveAssociations:                     false,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableNestedTransaction:                 true,
		AllowGlobalUpdate:                        true,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}

	db, err := gorm.Open(d, &dbConfig)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	Db = db.Debug()

	if conf.Config.Debug {
		Db.Logger = Db.Logger.LogMode(logger.Info)
	} else {
		Db.Logger = Db.Logger.LogMode(logger.Silent)
	}

	log.Infof("auto migrate tables")
	err = Db.AutoMigrate(
		&domain.CrawlerConf{},
		&domain.ProxyNode{},
		&domain.TohruFeed{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
