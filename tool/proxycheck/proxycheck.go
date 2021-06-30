package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/domain"
	"github.com/Luoxin/Eutamias/proxycheck"
	"github.com/Luoxin/Eutamias/proxynode"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/alexflint/go-arg"
	"github.com/gookit/color"
	"github.com/k0kubun/go-ansi"
	"github.com/olekukonko/tablewriter"
	"github.com/panjf2000/ants/v2"
	"github.com/schollz/progressbar/v3"
)

var cmdArgs struct {
	ConfigPath   string `arg:"-c,--config" help:"config file path"`
	FasterSpeed  bool   `arg:"-f,--fasterspeed" help:"order by speed"`
	LowerLatency bool   `arg:"-l,--lowerlatency" help:"order by delay"`
}

func main() {
	defer ants.Release()

	arg.MustParse(&cmdArgs)

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

	bar := progressbar.NewOptions(len(nodeList),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription("[cyan]proxy detection...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	checkDelay := proxycheck.NewProxyCheck()
	checkDelay.SetTimeout(time.Second * 3)
	checkDelay.SetCheckUrl("https://www.google.com")

	checkSpeed := proxycheck.NewProxyCheck()
	checkSpeed.SetTimeout(time.Second * 3)
	checkSpeed.SetCheckUrl("http://cachefly.cachefly.net/1mb.test")

	start := time.Now()

	var checkResultList CheckResultList
	var lock sync.Mutex
	var w sync.WaitGroup
	checkOnce := func(proxyNode *domain.ProxyNode) {
		defer w.Done()
		defer bar.Add(1)

		if proxyNode.DeathCount > 10 {
			return
		}

		if proxyNode.ProxySpeed < 0 {
			return
		}

		if proxyNode.ProxyNetworkDelay < 0 {
			return
		}

		var err error
		result := CheckResult{
			NodeName: utils.ShortStr(proxyNode.UrlFeature, 12),
		}

		result.Delay, _, err = checkDelay.CheckWithLink(proxyNode.Url)
		if err != nil {
			result.Delay = -1
		}

		_, result.Speed, err = checkSpeed.CheckWithLink(proxyNode.Url)
		if err != nil {
			result.Speed = -1
		}

		lock.Lock()
		checkResultList = append(checkResultList, &result)
		lock.Unlock()
	}

	pool, err := ants.NewPoolWithFunc(5, func(i interface{}) {
		checkOnce(i.(*domain.ProxyNode))
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
	w.Wait()

	checkResultList = checkResultList.Filter(func(result *CheckResult) bool {
		return result.Speed >= 0 && result.Delay >= 0
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
