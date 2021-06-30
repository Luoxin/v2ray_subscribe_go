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
	"github.com/schollz/progressbar/v3"
)

var cmdArgs struct {
	ConfigPath   string `arg:"-c,--config" help:"config file path"`
	FasterSpeed  bool   `arg:"-f,--fasterspeed" help:"order by speed"`
	LowerLatency bool   `arg:"-l,--lowerlatency" help:"order by delay"`
}

func main() {
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

	nodeList, err := proxynode.GetUsableNodeList(50, false, 1)
	if err != nil {
		color.Red.Printf("err:%v", err)
		return
	}

	bar := progressbar.NewOptions(len(nodeList),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[cyan][1/3][reset] proxy detection..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	check := proxycheck.NewProxyCheck()
	check.SetTimeout(time.Second * 5)

	var checkResultList CheckResultList
	var lock sync.Mutex
	var w sync.WaitGroup
	checkOnce := func(proxyNode *domain.ProxyNode) {
		defer w.Done()
		defer bar.Add(1)
		delay, speed, err := check.CheckWithLink(proxyNode.Url)
		if err != nil {
			lock.Lock()
			checkResultList = append(checkResultList, &CheckResult{
				NodeName: utils.ShortStr(proxyNode.UrlFeature, 12),
				Speed:    -1,
				Delay:    -1,
			})
			lock.Unlock()
		} else {
			lock.Lock()
			checkResultList = append(checkResultList, &CheckResult{
				NodeName: utils.ShortStr(proxyNode.UrlFeature, 12),
				Speed:    speed,
				Delay:    delay,
			})
			lock.Unlock()
		}
	}

	nodeList.Each(func(proxyNode *domain.ProxyNode) {
		w.Add(1)
		go checkOnce(proxyNode)
	})
	w.Wait()
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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"节点名",
		"速度",
		"时延",
	})

	for _, x := range checkResultList {
		table.Append([]string{
			x.NodeName,
			fmt.Sprintf("%.2f Kb/s", x.Speed),
			fmt.Sprintf("%.2f ms", x.Delay),
		})
	}
	table.Render()
}
