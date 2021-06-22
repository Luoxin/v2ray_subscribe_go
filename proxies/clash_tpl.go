package proxies

import (
	"path/filepath"
	"sync"

	"github.com/Luoxin/Eutamias/utils"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

type clashTplVal struct {
	Dns string
	ProxyList, ProxyNameList,
	CountryGroupList,
	NetEaseProxyList, NetEaseProxyNameList []string

	CountryNodeList []*countryNode

	TestUrl   string
	MixedPort uint32
}

var _lock sync.RWMutex

// TODO file watch
func init() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create {
					clashTplFile := filepath.Join(utils.GetExecPath(), "./resource/clashTpl")
					if utils.FileExists(clashTplFile) {
						buf, err := utils.FileRead(clashTplFile)
						if err == nil {
							log.Info("clash tpl changed")
							_lock.Lock()
							clashTpl = buf
							_lock.Unlock()
						}
					}
				}
			}
		}
	}()

	err = watcher.Add("./resource/clashTpl")
	if err != nil {
	    log.Errorf("err:%v", err)
	    return
	}

}

var clashTpl = `
mixed-port: {{ .MixedPort }}
allow-lan: true
mode: Rule
ipv6: true
experimental:
  ignore-resolve-fail: true
log-level: debug
external-controller: 127.0.0.1:9090
{{.Dns}}
proxies:
  - {"name":"网易音乐解锁","type":"http","server":"music.lolico.me","port":39000}
{{ range .NetEaseProxyList}}  - {{ .}}
{{ end}}{{ range .ProxyList}}  - {{ .}}
{{ end}}proxy-groups:
  - name: 🚀 节点选择
    type: select
    proxies:
      - 🔯 故障转移
      - ♻️ 自动选择
      - 🔮 负载均衡
      - 🚀 手动切换
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - DIRECT
  - name: 🚀 手动切换
    type: select
    proxies:
{{ range .ProxyNameList}}      - {{ .}}
{{ end}}  - name: ♻️ 自动选择
    type: url-test
    url: {{.TestUrl}}
    interval: 300
    tolerance: 50
    proxies:
{{ range .ProxyNameList}}      - {{ .}}
{{ end}}  - name: 🔯 故障转移
    type: fallback
    url: {{.TestUrl}}
    interval: 300
    tolerance: 50
    proxies:
{{ range .ProxyNameList}}      - {{ .}}
{{ end}}  - name: 🔮 负载均衡
    type: load-balance
    url: {{.TestUrl}}
    interval: 300
    tolerance: 50
    proxies:
{{ range .ProxyNameList}}      - {{ .}}
{{ end}}  - name: 📲 电报消息
    type: select
    proxies:
      - 🔯 故障转移
      - 🔮 负载均衡
      - ♻️ 自动选择
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
      - DIRECT
  - name: 📹 油管视频
    type: select
    proxies:
      - 🔯 故障转移
      - 🔮 负载均衡
      - ♻️ 自动选择
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
      - DIRECT
  - name: 🎥 奈飞视频
    type: select
    proxies:
      - 🎥 奈飞节点
      - 🚀 节点选择
      - ♻️ 自动选择
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
      - DIRECT
  - name: 📺 巴哈姆特
    type: fallback
    url: https://www.gamer.com.tw/
    interval: 300
    tolerance: 50
    proxies:
      - 🇹🇼 台湾省
      - 🚀 节点选择
      - 🚀 手动切换
      - DIRECT
  - name: 📺 哔哩哔哩
    type: fallback
    url: https://www.bilibili.com/
    interval: 300
    tolerance: 50
    proxies:
      - DIRECT
      - 🇹🇼 台湾省
  - name: 🌍 国外媒体
    type: select
    proxies:
      - 🚀 节点选择
      - ♻️ 自动选择
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
      - DIRECT
  - name: 🌏 国内媒体
    type: select
    proxies:
      - DIRECT
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
  - name: 📢 谷歌FCM
    type: select
    proxies:
      - 🔯 故障转移
      - ♻️ 自动选择
      - 🚀 节点选择
      - DIRECT
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
  - name: Ⓜ️ 微软云盘
    type: select
    proxies:
      - DIRECT
      - 🔯 故障转移
      - 🔮 负载均衡
      - ♻️ 自动选择
      - 🚀 节点选择
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
  - name: 🍎 苹果服务
    type: select
    proxies:
      - DIRECT
      - 🚀 节点选择
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
  - name: 🎮 游戏平台
    type: select
    proxies:
      - 🔯 故障转移
      - ♻️ 自动选择
      - 🚀 节点选择
      - DIRECT
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
  - name: 🎶 网易音乐
    type: fallback
    url: http://interface.music.163.com
    interval: 300
    tolerance: 50
    proxies:
      - 网易音乐解锁
{{ range .NetEaseProxyNameList}}      - {{ .}}
{{ end}}      - DIRECT
  - name: 🎯 全球直连
    type: select
    proxies:
      - 🔯 故障转移
      - ♻️ 自动选择
      - 🚀 节点选择
      - DIRECT
  - name: 🛑 广告拦截
    type: select
    proxies:
      - REJECT
      - DIRECT
  - name: 🍃 应用净化
    type: select
    proxies:
      - REJECT
      - DIRECT
  - name: 🆎 AdBlock
    type: select
    proxies:
      - REJECT
      - DIRECT
  - name: 🛡️ 隐私防护
    type: select
    proxies:
      - REJECT
      - DIRECT
  - name: 🐟 漏网之鱼
    type: load-balance
    url: {{.TestUrl}}
    interval: 300
    tolerance: 50
    proxies:
      - 🚀 节点选择
      - 🔯 故障转移
      - ♻️ 自动选择
      - DIRECT
{{ range .CountryGroupList}}      - {{ .}}
{{ end}}      - 🚀 手动切换
{{ range .CountryNodeList}}  - name: {{.Emoji}} {{.Name}}
    type: fallback
    url: {{.TestUrl}}
    interval: 300
    tolerance: 50
    proxies:
{{ range .NameList}}      - {{.}}
{{ end}}{{ end}}  - name: 🎥 奈飞节点
    type: select
    proxies:
      - DIRECT
rule-providers:
  reject:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/reject.txt"
    path: ./ruleset/reject.yaml
    interval: 86400

  icloud:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/icloud.txt"
    path: ./ruleset/icloud.yaml
    interval: 86400

  apple:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/apple.txt"
    path: ./ruleset/apple.yaml
    interval: 86400

  google:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/google.txt"
    path: ./ruleset/google.yaml
    interval: 86400

  proxy:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/proxy.txt"
    path: ./ruleset/proxy.yaml
    interval: 86400

  direct:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/direct.txt"
    path: ./ruleset/direct.yaml
    interval: 86400

  private:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/private.txt"
    path: ./ruleset/private.yaml
    interval: 86400

  gfw:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/gfw.txt"
    path: ./ruleset/gfw.yaml
    interval: 86400

  greatfire:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/greatfire.txt"
    path: ./ruleset/greatfire.yaml
    interval: 86400

  tld-not-cn:
    type: http
    behavior: domain
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/tld-not-cn.txt"
    path: ./ruleset/tld-not-cn.yaml
    interval: 86400

  telegramcidr:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/telegramcidr.txt"
    path: ./ruleset/telegramcidr.yaml
    interval: 86400

  cncidr:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/cncidr.txt"
    path: ./ruleset/cncidr.yaml
    interval: 86400

  lancidr:
    type: http
    behavior: ipcidr
    url: "https://cdn.jsdelivr.net/gh/Loyalsoldier/clash-rules@release/lancidr.txt"
    path: ./ruleset/lancidr.yaml
    interval: 86400
rules:
  - GEOIP,CN,DIRECT
  - MATCH,🐟 漏网之鱼
`
