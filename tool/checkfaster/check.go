package main

import (
	"fmt"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	"github.com/d-tsuji/clipboard"

	"subsrcibe/parser"
	"subsrcibe/proxycheck"
)

func main() {
	logFormatter := &nested.Formatter{
		FieldsOrder: []string{
			log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
			log.FieldKeyFunc, log.FieldKeyMsg,
		},
		TimestampFormat:  time.RFC3339,
		HideKeys:         true,
		NoFieldsSpace:    true,
		NoUppercaseLevel: true,
		TrimMessages:     true,
		CallerFirst:      true,
	}

	log.SetFormatter(logFormatter)
	log.SetReportCaller(true)

	err := func() error {
		resp, err := nic.Get("http://127.0.0.0:8080", nil)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		check := proxycheck.NewProxyCheck()
		err = check.Init()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = check.SetMaxSize(100)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		var resultList proxycheck.ResultList
		parser.NewFuzzyMatchingParser().ParserText(resp.Text).Unique().Each(func(s string) {
			err := check.Add(s, func(result proxycheck.Result) error {
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
			fmt.Println(result.ProxyUrl)
			text += fmt.Sprintf("%s\n", result.ProxyUrl)
		})

		err = clipboard.Set(text)
		if err != nil {
			log.Errorf("err:%v", err)
		}

		return nil
	}()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
