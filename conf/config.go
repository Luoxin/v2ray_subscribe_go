package conf

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Luoxin/Eutamias/utils"
	"github.com/elliotchance/pie/pie"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var roleList = pie.Strings{
	"Kobayashi-san",
	"Tohru",
}

type base struct {
	Role                string `yaml:"role" json:"role"`
	KobayashiSanAddr    string `yaml:"Kobayashi-san_addr" json:"Kobayashi-san_addr"`
	KobayashiSanHomeKey string `yaml:"Kobayashi-san_home_key" json:"Kobayashi-san_home_key"`
	TohruKey            string `yaml:"tohru_key" json:"tohru_key"`
	TohruPassword       string `yaml:"tohru_password" json:"tohru_password"`
}

type db struct {
	Typ      string `yaml:"type" json:"type"`
	Host     string `yaml:"host" json:"host"`
	Port     uint32 `yaml:"port" json:"port"`
	Database string `yaml:"database" json:"database"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
}

type crawler struct {
	Enable  bool   `yaml:"enable" json:"enable"`
	Proxies string `yaml:"proxies" json:"proxies"`

	CrawlerInterval uint32 `yaml:"crawler_interval" json:"crawler_interval"`
}

type proxyCheck struct {
	Enable        bool   `yaml:"enable" json:"enable"`
	CheckInterval uint32 `yaml:"check_interval" json:"check_interval"`
}

type httpService struct {
	Enable bool   `yaml:"enable" json:"enable"`
	Host   string `yaml:"host" json:"host"`
	Port   uint32 `yaml:"port" json:"port"`
}

type proxy struct {
	Enable    bool   `yaml:"enable" json:"enable"`
	MixedPort uint32 `yaml:"mixed-port" json:"mixed-port"`
}

type dns struct {
	Enable        bool        `yaml:"enable" json:"enable"`
	EnableService bool        `yaml:"enable_service" json:"enable_service"`
	ServicePort   uint32      `yaml:"service_port" json:"service_port"`
	Nameserver    pie.Strings `yaml:"nameserver" json:"nameserver"`
}

type cache struct {
	Typ string `yaml:"type" json:"type"`
}

type config struct {
	Base base `yaml:"base" json:"base"`

	Debug bool `yaml:"debug" json:"debug"`

	Db    db    `yaml:"db" json:"db"`
	Cache cache `yaml:"cache" json:"cache"`

	Crawler     crawler     `yaml:"crawler" json:"crawler"`
	ProxyCheck  proxyCheck  `yaml:"proxy_check" json:"proxy_check"`
	HttpService httpService `yaml:"http_service" json:"http_service"`

	Proxy proxy `yaml:"proxy" json:"proxy"`

	Dns dns `yaml:"dns" json:"dns"`
}

func (p config) IsTohru() bool {
	return p.Base.Role == "Tohru"
}

func (p config) KobayashiSan() bool {
	return p.Base.Role == "Kobayashi-san"
}

var Config config

func InitConfig(configFilePatch string) error {
	execPath := utils.GetExecPath()

	if configFilePatch == "" {
		// 可能存在的目录
		viper.AddConfigPath("./")
		viper.AddConfigPath("../")
		viper.AddConfigPath("./conf/")
		viper.AddConfigPath("../conf/")
		{
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Errorf("err:%v", err)
			} else {
				viper.AddConfigPath(filepath.Join(homeDir, "Eutamias"))
			}
		}

		viper.AddConfigPath(execPath)

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	} else {
		log.Info(configFilePatch)
		viper.SetConfigFile(configFilePatch)
	}

	// 配置一些默认值
	viper.SetDefault("base.role", "Kobayashi-san")
	viper.SetDefault("base.Kobayashi-san_addr", "http://0.0.0.0:8080")
	viper.SetDefault("base.Kobayashi-san_home_key", "T6Z14ey@rj)?LjMvkih+?.W}JAU?V{qvsD+H_)R/")

	viper.SetDefault("http_service.enable", true)
	viper.SetDefault("http_service.port", 8080)
	viper.SetDefault("http_service.host", "127.0.0.1")

	viper.SetDefault("debug", false)

	viper.SetDefault("db.type", "sqlite")

	viper.SetDefault("db.cache", "")

	viper.SetDefault("crawler.enable", true)
	viper.SetDefault("crawler.proxies", "http://127.0.0.1:7890")
	viper.SetDefault("crawler.crawler_interval", 3600)

	viper.SetDefault("proxy_check.enable", true)
	viper.SetDefault("proxy_check.check_interval", 300)

	viper.SetDefault("proxy.enable", false)
	viper.SetDefault("proxy.mixed-port", 7890)

	viper.SetDefault("profiler.enable", false)
	viper.SetDefault("profiler.server_address", "http://localhost:4040")

	viper.SetDefault("dns.enable", false)
	viper.SetDefault("dns.enable_service", false)
	viper.SetDefault("dns.service_port", 53)
	viper.SetDefault("dns.nameserver", []string{"1.2.4.8", "223.5.5.5", "8.8.8.8", "tls://dns.alidns.com"})

	err := viper.ReadInConfig()
	if err != nil {
		switch e := err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warnf("not found conf file, use default")
		case *os.PathError:
			log.Warnf("not find conf file in %s", e.Path)
		default:
			log.Errorf("err:%v", err)
			return err
		}

	}

	j, err := json.Marshal(viper.AllSettings())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = json.Unmarshal(j, &Config)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = viper.WriteConfigAs(filepath.Join(execPath, "config.yaml"))
	if err != nil {
		log.Errorf("err:%v", err)
	}

	if !roleList.Contains(Config.Base.Role) {
		Config.Base.Role = "Kobayashi-san"
	}

	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Debugf("get config %+v", Config)

	err = Ecc.Init(Config.Base.KobayashiSanHomeKey)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = os.MkdirAll(utils.GetConfigDir(), 0777)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
