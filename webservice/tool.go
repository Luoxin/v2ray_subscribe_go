package webservice

import (
	// "github.com/gofiber/storage/mysql"

	"github.com/gofiber/fiber/v2"
)

func InitStorage(dbAddr string) (fiber.Storage, error) {
	// addrList := strings.Split(dbAddr, "://")
	// if len(addrList) < 2 {
	// 	log.Errorf("Wrong database address")
	// 	return nil, errors.New("invalid args")
	// }
	//
	var storage fiber.Storage
	// switch strings.ToLower(addrList[0]) {
	// case "sqlite":
	// 	// d = sqlite.Open(strings.Join(addrList[1:], ""))
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

	return storage, nil
}
