package main

import (
	"crypto/sha512"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/crawler"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/parser"
	"github.com/Luoxin/Eutamias/proxycheck"
	"github.com/Luoxin/Eutamias/proxynode"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/alexflint/go-arg"
	"github.com/gookit/color"
	log "github.com/sirupsen/logrus"

	// "github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
	"github.com/panjf2000/ants/v2"
	"github.com/pterm/pterm"
	// "github.com/schollz/progressbar/v3"
)

var cmdArgsstruct {
	SubUrl       string `arg:"-u,--suburl" help:"sub url"`
	UseProxy     bool   `arg:"-p,--useproxy" help:"sub use proxy"`
	ConfigPath   string `arg:"-c,--config" help:"config file path"`
	FasterSpeed  bool   `arg:"-f,--fasterspeed" help:"order by speed"`
	LowerLatency bool   `arg:"-l,--lowerlatency" help:"order by delay"`
	Debug        bool   `arg:"-d,--debug" help:"enable debug"`
}

func main() {
	log.SetReportCaller(false)
	defer ants.Release()
	start := time.Now()

	arg.MustParse(&cmdArgs)
	if cmdArgs.Debug {
		log.SetReportCaller(true)
	}

	checkDelay := proxycheck.NewProxyCheck()
	checkDelay.SetTimeout(time.Second)
	checkDelay.SetCheckUrl("https://www.google.com")

	checkSpeed := proxycheck.NewProxyCheck()
	checkSpeed.SetTimeout(time.Second)
	checkSpeed.SetCheckUrl("http://cachefly.cachefly.net/1mb.test")

	var lock sync.Mutex

	var checkResultList CheckResultList
	var w sync.WaitGroup

	// bar := progressbar.NewOptions(0,
	// 	progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
	// 	progressbar.OptionEnableColorCodes(true),
	// 	progressbar.OptionShowBytes(true),
	// 	progressbar.OptionSetWidth(50),
	// 	progressbar.OptionSetDescription("[cyan]proxy detection...[reset]"),
	// 	progressbar.OptionSetTheme(progressbar.Theme{
	// 		Saucer:        "[green]=[green]",
	// 		SaucerHead:    "[green]>[yellow]",
	// 		SaucerPadding: "[yellow]=[reset]",
	// 		BarStart:      "[",
	// 		BarEnd:        "]",
	// 	}),
	// 	progressbar.OptionSetVisibility(true),
	// 	progressbar.OptionEnableColorCodes(true),
	// 	progressbar.OptionClearOnFinish(),
	// 	progressbar.OptionShowIts(),
	// 	progressbar.OptionSetItsString("nodes"),
	// 	progressbar.OptionThrottle(time.Second),
	// 	progressbar.OptionFullWidth(),
	// )

	var p *pterm.ProgressbarPrinter
	var err error

	initPterm := func(total int) {
		p, err = pterm.DefaultProgressbar.
			WithTitle("节点检测").
			WithShowElapsedTime(true).
			WithShowTitle(true).
			WithShowCount(true).
			WithTotal(total).
			Start()
		if err != nil {
			log.Fatalf("err:%v", err)
			return
		}
	}

	checkOnce := func(nodeUrl string) {
		start := time.Now()
		defer w.Done()
		// defer bar.Add(1)
		defer p.Increment()

		var err error
		result := CheckResult{
			NodeName: utils.ShortStr(fmt.Sprintf("%x", sha512.Sum512([]byte(nodeUrl))), 12),
		}

		p.Title = fmt.Sprintf("check %v", result.NodeName)

		result.Delay, _, err = checkDelay.CheckWithLink(nodeUrl)
		if err != nil {
			result.Delay = -1
		}

		_, result.Speed, err = checkSpeed.CheckWithLink(nodeUrl)
		if err != nil {
			result.Speed = -1
		}

		lock.Lock()
		checkResultList = append(checkResultList, &result)
		if result.Speed >= 0 && result.Delay >= 0 {
			pterm.Success.Printfln("节点%v(检测耗时:%v),速度:%.3fKb/s,延时:%.3fms",
				result.NodeName, time.Since(start), result.Speed, result.Delay)
		} else if result.Speed < 0 && result.Delay < 0 {
			pterm.Error.Printfln("节点%v(检测耗时:%v)无法联通", result.NodeName, time.Since(start))
		} else {
			pterm.Warning.Printfln("节点%v(检测耗时:%v)%v%v",
				result.NodeName, time.Since(start),
				func() string {
					if result.Speed < 0 {
						return ""
					}
					return fmt.Sprintf(",速度:%.3fKb/s", result.Speed)
				}(), func() string {
					if result.Delay < 0 {
						return ""
					}
					return fmt.Sprintf(",延时:%.3fms", result.Delay)
				}())
		}
		lock.Unlock()
	}

	switch cmdArgs.SubUrl {
	case "":
		err := conf.InitConfig(cmdArgs.ConfigPath)
		if err != nil {
			color.Red.Printf("err:%v", err)
			return
		}

		err = db.InitDb()
		if err != nil {
			color.Red.Printf("err:%v", err)
			return
		}

		nodeList, err := proxynode.GetUsableNodeList(100, false, 1)
		if err != nil {
			color.Red.Printf("err:%v", err)
			return
		}

		// bar.ChangeMax(len(nodeList))
		initPterm(len(nodeList))

		pool, err := ants.NewPoolWithFunc(20, func(i interface{}) {
			checkOnce(i.(*domain.ProxyNode).Url)
		})
		if err != nil {
			color.Red.Printf("err:%v", err)
			return
		}

		nodeList.Each(func(proxyNode *domain.ProxyNode) {
			w.Add(1)
			err = pool.Invoke(proxyNode)
			if err != nil {
				color.Red.Printf("err:%v", err)
				w.Done()
			}
		})
	default:
		conf.Config.Crawler.Proxies = "http://127.0.0.1:7890"
		rspBody, err := crawler.NewHttpDownloader().
			Download("GET", cmdArgs.SubUrl, nil, domain.CrawlerConf_Rule{
				UseProxy: cmdArgs.UseProxy,
			})
		if err != nil {
			color.Red.Printf("err:%v", err)
			return
		}

		nodeList := parser.NewFuzzyMatchingParser().
			ParserText(rspBody).
			Filter(func(s string) bool {
				return strings.Contains(s, "://")
			})

		// bar.ChangeMax(len(nodeList))
		initPterm(len(nodeList))

		pool, err := ants.NewPoolWithFunc(20, func(i interface{}) {
			checkOnce(i.(string))
		})
		if err != nil {
			color.Red.Printf("err:%v", err)
			return
		}

		nodeList.Each(func(nodeUrl string) {
			w.Add(1)
			err = pool.Invoke(nodeUrl)
			if err != nil {
				color.Red.Printf("err:%v", err)
				w.Done()
			}
		})
	}

	w.Wait()
	// bar.Close()

	total := len(checkResultList)
	checkResultList = checkResultList.Filter(func(result *CheckResult) bool {
		return result.Speed >= 0 || result.Delay >= 0
	})

	if cmdArgs.FasterSpeed {
		checkResultList = checkResultList.SortUsing(func(a, b *CheckResult) bool {
			return a.Speed > b.Speed
		})
	}

	if cmdArgs.LowerLatency {
		checkResultList = checkResultList.SortUsing(func(a, b *CheckResult) bool {
			return a.Delay < b.Delay
		})
	}

	fmt.Println("")

	fmt.Printf("used:%v\n", time.Since(start))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetCaption(true, fmt.Sprintf("测速结果(可用率：%.2f%%)", float64(len(checkResultList))/float64(total)))
	table.SetHeader([]string{
		"节点名",
		"速度",
		"时延",
	})

	for _, x := range checkResultList {
		table.Append([]string{
			x.NodeName,
			func() string {
				if x.Speed < 0 {
					return "-"
				}
				return fmt.Sprintf("%.2f Kb/s", x.Speed)
			}(),
			func() string {
				if x.Delay < 0 {
					return "-"
				}
				return fmt.Sprintf("%.2f ms", x.Delay)
			}(),
		})
	}
	table.Render()
}
