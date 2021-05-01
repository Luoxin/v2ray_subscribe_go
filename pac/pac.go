package pac

import (
	"bytes"
	"encoding/base64"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/elliotchance/pie/pie"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/utils"
)

type pac struct {
	js       string
	updateAt uint32
	client   *resty.Client

	writeLock  sync.RWMutex
	updateLock sync.Mutex
}

var Pac = NewPac()

func InitPac() {
	Pac = NewPac()
	Pac.UpdatePac()
}

func NewPac() *pac {
	return &pac{
		client: resty.New().
			SetRetryCount(2).
			SetTimeout(time.Second * 5).
			SetRetryMaxWaitTime(time.Second * 5).
			SetRetryWaitTime(time.Second).
			SetLogger(log.New()),
	}
}

func (p *pac) Get() string {
	if !p.needUpdate() {
		p.writeLock.RLock()
		defer p.writeLock.RUnlock()
		return p.js
	}

	p.UpdatePac()
	return p.js
}

func (p *pac) write(js string) {
	p.writeLock.Lock()
	defer p.writeLock.Unlock()
	if js == "" {
		return
	}

	p.js = js
	p.updateAt = utils.Now()
}

func (p *pac) needUpdate() bool {
	return p.js == "" || utils.Now()-p.updateAt > 86400
}

func (p *pac) UpdatePac() {
	p.updateLock.Lock()
	defer p.updateLock.Unlock()

	if !p.needUpdate() {
		return
	}

	js := p.buildPac("PROXY 127.0.0.1:7890;", "DIRECT_PROXY", p.getRuleList())
	p.write(js)
	go func() {
		p.update2Sys()
	}()
}

func (p *pac) update2Sys() {
	// // 没有开启http service，也就意味着没有pac
	// if !conf.Config.HttpService.Enable {
	// 	return
	// }
	//
	// if runtime.GOOS == "windows" {
	// 	pacHost := conf.Config.HttpService.Host
	// 	if pacHost == "0.0.0.0" {
	// 		pacHost = "127.0.0.1"
	// 	}
	//
	// 	utils.SetProxy(fmt.Sprintf("127.0.0.1:%d", conf.Config.Proxy.MixedPort), fmt.Sprintf("http://%s:%d/api/subscribe/pac?_=%d", pacHost, conf.Config.HttpService.Port, p.updateAt))
	// 	log.Infof("set system PAC finish(%d)", p.updateAt)
	// }
}

func (p *pac) buildPac(proxy, defaultWay string, ruleList pie.Strings) string {
	t, err := template.New("").Parse(pacTpl)
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}

	var b bytes.Buffer
	err = t.Execute(&b, map[string]interface{}{
		"RuleList": ruleList,
		"Proxy":    proxy,
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return ""
	}
	return b.String()
}

func (p *pac) getRuleList() pie.Strings {
	var gfwUrlList = pie.Strings{
		"https://cdn.jsdelivr.net/gh/gfwlist/gfwlist@master/gfwlist.txt",
		"https://pagure.io/gfwlist/raw/master/f/gfwlist.txt",
		"https://repo.or.cz/gfwlist.git/blob_plain/HEAD:/gfwlist.txt",
		"https://bitbucket.org/gfwlist/gfwlist/raw/HEAD/gfwlist.txt",
		"https://git.tuxfamily.org/gfwlist/gfwlist.git/plain/gfwlist.txt",
	}

	ruleMap := map[string]bool{}
	gfwUrlList.Each(func(s string) {
		log.Debugf("handle pac for:%v", s)
		resp, err := p.client.R().Get(s)
		if err != nil {
			return
		}

		ruleList := strings.Split(resp.String(), "\n")

		pie.Strings(ruleList).Each(func(str string) {
			lineByte, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				log.Errorf("err:%v", err)
				return
			}

			line := string(lineByte)
			if !strings.HasPrefix(line, "!") {
				ruleMap[line] = true
			}
		})
	})

	var ruleList []string
	for rule := range ruleMap {
		ruleList = append(ruleList, rule)
	}

	return ruleList
}

func Get() string {
	return Pac.Get()
}
