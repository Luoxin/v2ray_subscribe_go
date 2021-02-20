package geolite

import (
	"errors"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
	log "github.com/sirupsen/logrus"

	country2 "subscribe/country"
)

var db *geoip2.Reader

func init() {
	var err error
	db, err = geoip2.Open("./GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatalf("err:%v", err)
	}
}

type Country struct {
	Code   string
	CnName string
	EnName string
	Emoji  string
}

func GetCountry(host string) (*Country, error) {
	ns, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}

	if len(ns) == 0 {
		return nil, errors.New("not found ip")
	}

	record, err := db.Country(net.ParseIP(ns[0].String()))
	if err != nil {
		return nil, err
	}

	c := Country{
		CnName: record.Country.Names["zh-CN"],
		EnName: record.Country.Names["en"],
	}

	country := country2.Country.GetByEnName(strings.ReplaceAll(record.Country.Names["en"], " ", ""))
	if country == nil {
		country = country2.Country.GetByCnName(record.Country.Names["zh-CN"])
	}

	if country != nil {
		c.Code = country.Code
		c.Emoji = country.Emoji
		c.CnName = country.CnName
		c.EnName = country.EnName
	} else {
		log.Warnf("not found country:%+v", record)
	}

	return &c, nil
}
