package geolite

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"
)

var db *geoip2.Reader

func init() {
	var err error
	db, err = geoip2.Open("./GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatalf("err:%v", err)
	}
}

func GetCountry(host string) *geoip2.Country {
	ns, err := net.LookupIP(host)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil
	}

	if len(ns) == 0 {
		return nil
	}

	record, err := db.Country(net.ParseIP(ns[0].String()))
	if err != nil {
		log.Fatal(err)
	}

	return record
}
