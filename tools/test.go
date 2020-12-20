package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

func main() {
	db, err := geoip2.Open("D:\\Coding\\v2ray_subscribe_go\\tools\\GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//fmt.Println(db.Metadata())

	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP("1.208.0.0")
	record, err := db.Country(ip)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(record.Country.Names)
}