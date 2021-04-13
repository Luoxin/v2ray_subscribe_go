package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
	"github.com/d-tsuji/clipboard"
	"github.com/eddieivan01/nic"
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/node"
	"github.com/Luoxin/Eutamias/parser"
	"github.com/Luoxin/Eutamias/proxycheck"
)

func main() {
	var cmdArgs struct {
		SubUrl string `arg:"-u" help:"sub url"`
	}

	arg.MustParse(&cmdArgs)

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
		if cmdArgs.SubUrl != "" {
			resp, err := nic.Get(cmdArgs.SubUrl, nil)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			nodeUrl = parser.NewFuzzyMatchingParser().ParserText(resp.Text).Unique()
		} else {
			err := conf.InitConfig("")
			if err != nil {
				log.Fatalf("init config err:%v", err)
				return err
			}

			err = db.InitDb(conf.Config.Db.Addr)
			if err != nil {
				log.Fatalf("init db err:%v", err)
				return err
			}

			nodes, err := node.GetUsableNodeList(50, true)
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}

			nodes.Each(func(node *domain.ProxyNode) {
				nodeUrl = append(nodeUrl, node.Url)
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
