package db

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/Luoxin/Eutamias/utils"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/domain"
)

var Db *gorm.DB

func InitDb() error {
	// log.Infof("connect to database %v", dbAddr)
	dbConfig := conf.Config.Db

	var d gorm.Dialector
	switch dbConfig.Typ {
	case "sqlite":
		d = sqlite.Open(fmt.Sprintf("%s?check_same_thread=false", filepath.ToSlash(filepath.Join(utils.GetExecPath(), ".eutamias.es"))))
	case "mysql":
		dsn := fmt.Sprintf("%s:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local&checkConnLiveness=true&writeTimeout=3s&timeout=5s&readTimeout=30s",
			dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
		log.Infof("connect to database %v", dsn)
		d = mysql.New(mysql.Config{
			DSN:                    dsn,
			DefaultStringSize:      256,
			DontSupportRenameIndex: true,
		})
	default:
		return errors.New("database types are not supported")
	}

	// addrList := strings.Split(dbAddr, "://")
	// if len(addrList) < 2 {
	// 	log.Errorf("Wrong database address")
	// 	return errors.New("invalid args")
	// }
	//

	// switch strings.ToLower(addrList[0]) {
	// case "sqlite":

	// case "mysql":

	// default:
	// 	return errors.New("unsupported database")
	// }

	db, err := gorm.Open(d, &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: &schema.NamingStrategy{
			TablePrefix:   "eutamias_",
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
	})
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
