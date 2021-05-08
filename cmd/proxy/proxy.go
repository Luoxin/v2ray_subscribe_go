package main

import (
	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/geolite"
	"github.com/Luoxin/Eutamias/node"
	"github.com/Luoxin/Eutamias/proxy"
	"github.com/Luoxin/Eutamias/utils"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := conf.InitConfig("")
	if err != nil {
		log.Fatalf("init config err:%v", err)
		return
	}

	log.Info("init conf success")

	err = geolite.InitGeoLite()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	log.Info("init geolite success")

	err = db.InitDb(conf.Config.Db.Addr)
	if err != nil {
		log.Fatalf("init db err:%v", err)
		return
	}

	nodes, err := node.GetUsableNodeList(100, false, 1)
	if err != nil {
		log.Fatalf("init db err:%v", err)
		return
	}

	nodes.Each(func(proxyNode *domain.ProxyNode) {
		p, err := proxy.ParseProxy(proxyNode.Url)
		if err != nil {
			log.Fatalf("init db err:%v", err)
			return
		}

		log.Infof("%v %v", utils.ShortStr(proxyNode.UrlFeature, 12), p.BaseInfo().Server)
	})
}
