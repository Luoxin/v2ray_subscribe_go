package main

import (
	"flag"
	"fmt"

	"github.com/d-tsuji/clipboard"
	"github.com/eddieivan01/nic"
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"

	"subscribe/conf"
	"subscribe/db"
	"subscribe/domain"
	"subscribe/http"
	"subscribe/parser"
	"subscribe/proxycheck"
)

var subUrl string

func init() {
	flag.StringVar(&subUrl, "u", "", "web for maybe has url")
}

func main() {
	err := func() error {
		check := proxycheck.NewProxyCheck()
		err := check.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = check.SetMaxSize(100)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		var nodeUrl pie.Strings
		if subUrl != "" {
			resp, err := nic.Get(subUrl, nil)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			nodeUrl = parser.NewFuzzyMatchingParser().ParserText(resp.Text).Unique()
		} else {
			err := conf.InitConfig()
			if err != nil {
				log.Fatalf("init config err:%v", err)
				return err
			}

			err = db.InitDb(conf.Config.Db.Addr)
			if err != nil {
				log.Fatalf("init db err:%v", err)
				return err
			}

			nodes, err := http.GetUsableNodeList()
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			nodes.Each(func(node *domain.ProxyNode) {
				if node.NodeDetail == nil {
					return
				}

				nodeUrl = append(nodeUrl, node.NodeDetail.Buf)
			})
		}

		var resultList proxycheck.ResultList
		nodeUrl.Each(func(s string) {
			err := check.AddWithLink(s, func(result proxycheck.Result) error {
				resultList = append(resultList, &result)
				return nil
			})
			if err != nil {
				log.Errorf("err:%v", err)
			}
		})

		check.WaitFinish()

		var text string
		resultList.Filter(func(result *proxycheck.Result) bool {
			return result.Err == nil
		}).SortUsing(func(a, b *proxycheck.Result) bool {
			return a.Delay < b.Delay
		}).SortUsing(func(a, b *proxycheck.Result) bool {
			return a.Speed > b.Speed
		}).Top(5).Each(func(result *proxycheck.Result) {
			text += fmt.Sprintf("%s\n", result.ProxyUrl)
		})

		err = clipboard.Set(text)
		if err != nil {
			log.Errorf("err:%v", err)
		}

		fmt.Println(text)
		return nil
	}()
	if err != nil {
		log.Errorf("err:%v", err)
	}
}
