package webservice

import (
	"path/filepath"
	"time"

	"github.com/Luoxin/Eutamias/utils"
	"github.com/gofiber/storage/sqlite3"
)

func InitStorage(dbAddr string) error {
	// addrList := strings.Split(dbAddr, "://")
	// if len(addrList) < 2 {
	// 	log.Errorf("Wrong database address")
	// 	return nil, errors.New("invalid args")
	// }
	//

	// switch strings.ToLower(addrList[0]) {
	// case "sqlite":

	storage = sqlite3.New(sqlite3.Config{
		Database:   filepath.Join(utils.GetExecPath(), ".eutamias.es"),
		Table:      "fiber_storage",
		Reset:      false,
		GCInterval: time.Hour,
	})

	// case "mysql":
	// 	storage = mysql.New(mysql.Config{
	// 		Host:       "",
	// 		Port:       0,
	// 		Username:   "",
	// 		Password:   "",
	// 		Database:   "",
	// 		Table:      "subscribe_fiber",
	// 		Reset:      false,
	// 		GCInterval: 0,
	// 	})
	//
	// 	// d = mysql.Open(strings.Join(addrList[1:], ""))
	// default:
	// 	return nil, errors.New("unsupported database")
	// }

	return nil
}
