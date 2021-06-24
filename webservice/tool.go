package webservice

import (
	"path/filepath"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/gofiber/storage/memory"
	"github.com/gofiber/storage/mysql"
	"github.com/gofiber/storage/postgres"
	"github.com/gofiber/storage/sqlite3"
)

func InitStorage() error {
	dbConfig := conf.Config.Db

	switch dbConfig.Typ {
	case "sqlite":
		storage = sqlite3.New(sqlite3.Config{
			Database:   filepath.Join(utils.GetExecPath(), ".eutamias.es"),
			Table:      "eutamias_fiber_storage",
			Reset:      false,
			GCInterval: time.Hour,
		})

	case "mysql":
		storage = mysql.New(mysql.Config{
			Host:       dbConfig.Host,
			Port:       int(dbConfig.Port),
			Username:   dbConfig.User,
			Password:   dbConfig.Password,
			Database:   dbConfig.Database,
			Table:      "eutamias_fiber_storage",
			Reset:      false,
			GCInterval: time.Hour,
		})

	case "postgres":
		storage = postgres.New(postgres.Config{
			Host:       dbConfig.Host,
			Port:       int(dbConfig.Port),
			Username:   dbConfig.User,
			Password:   dbConfig.Password,
			Database:   dbConfig.Database,
			Table:      "eutamias_fiber_storage",
			Reset:      false,
			GCInterval: time.Hour,
			SslMode:    "disable",
		})

	default:
		storage = memory.New(memory.Config{
			GCInterval: time.Hour,
		})
	}

	return nil
}
