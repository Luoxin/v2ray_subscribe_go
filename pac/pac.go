package pac

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"github.com/eddieivan01/nic"
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"
	"io"
	"strings"
	"subscribe/utils"
	"sync"
	"text/template"
)

type pac struct {
	js       string
	updateAt uint32
}

var Pac = NewPac()
var lock sync.RWMutex

func InitPac() {
	Pac.UpdatePac()
}

func NewPac() *pac {
	return &pac{}
}

func (p *pac) Get() string {
	if utils.Now()-p.updateAt < 86400 {
		return p.js
	}

	p.UpdatePac()
	return p.js
}

func (p *pac) UpdatePac() {
	lock.Lock()
	lock.Unlock()

	js := p.buildPac("PROXY 127.0.0.1:7890;", "DIRECT_PROXY", p.getRuleList())
	if js != "" {
		p.js = js
		p.updateAt = utils.Now()
	}
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
		"http://cdn.jsdelivr.net/gh/gfwlist/gfwlist@master/gfwlist.txt",
		"https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt",
		"https://pagure.io/gfwlist/raw/master/f/gfwlist.txt",
		"http://repo.or.cz/gfwlist.git/blob_plain/HEAD:/gfwlist.txt",
		"https://bitbucket.org/gfwlist/gfwlist/raw/HEAD/gfwlist.txt",
		"https://gitlab.com/gfwlist/gfwlist/raw/master/gfwlist.txt",
		"https://git.tuxfamily.org/gfwlist/gfwlist.git/plain/gfwlist.txt",
	}

	var ruleList pie.Strings

	gfwUrlList.Each(func(s string) {
		resp, err := nic.Get(s, nil)
		if err != nil {
			return
		}

		decoder := base64.NewDecoder(base64.StdEncoding, resp.Body)
		reader := bufio.NewReader(decoder)
		for {

			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}

			l := string(line)

			if !strings.HasPrefix(l, "!") {
				ruleList = append(ruleList, l)
			}
		}
		resp.Body.Close()

	})

	return ruleList.Unique()
}

func Get() string {
	return Pac.Get()
}
