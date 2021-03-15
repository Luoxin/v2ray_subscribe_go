package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/elliotchance/pie/pie"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const Version = "0.0.0.4"

var roleList = pie.Strings{
	"Kobayashi-san",
	"Tohru",
}

type base struct {
	Role                string `yaml:"role" json:"role"`
	KobayashiSanAddr    string `yaml:"Kobayashi-san_addr" json:"Kobayashi-san_addr"`
	KobayashiSanHomeKey string `yaml:"Kobayashi-san_home_key" json:"Kobayashi-san_home_key"`
}

type db struct {
	Addr string `yaml:"addr" json:"addr"`
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

type _profiler struct {
	Enable        bool   `yaml:"enable" json:"enable"`
	ServerAddress string `yaml:"server_address" json:"server_address"`
}

type config struct {
	Base base `yaml:"base" json:"base"`

	Debug bool `yaml:"debug" json:"debug"`

	Db db `yaml:"db" json:"db"`

	Crawler     crawler     `yaml:"crawler" json:"crawler"`
	ProxyCheck  proxyCheck  `yaml:"proxy_check" json:"proxy_check"`
	HttpService httpService `yaml:"http_service" json:"http_service"`

	Proxy proxy `yaml:"proxy" json:"proxy"`

	Profiler _profiler `yaml:"profiler" json:"profiler"`
}

func (p config) IsTohru() bool {
	return p.Base.Role == "Tohru"
}

func (p config) KobayashiSan() bool {
	return p.Base.Role == "Kobayashi-san"
}

var Config config

var LogFormatter = &nested.Formatter{
	FieldsOrder: []string{
		log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
		log.FieldKeyFunc, log.FieldKeyMsg,
	},
	CustomCallerFormatter: func(f *runtime.Frame) string {
		return fmt.Sprintf("(%s %s:%d)", f.Function, path.Base(f.File), f.Line)
	},
	TimestampFormat:  time.RFC3339,
	HideKeys:         true,
	NoFieldsSpace:    true,
	NoUppercaseLevel: true,
	TrimMessages:     true,
	CallerFirst:      true,
}

func InitConfig() error {
	log.SetFormatter(LogFormatter)
	log.SetReportCaller(true)

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
			viper.AddConfigPath(filepath.Join(homeDir, "subscribe"))
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 配置一些默认值
	viper.SetDefault("base.role", "Kobayashi-san")
	viper.SetDefault("base.Kobayashi-san_addr", "http://127.0.0.1:8080/api/subscribe")
	viper.SetDefault("base.Kobayashi-san_home_key", "T6Z14ey@rj)?LjMvkih+?.W}JAU?V{qvsD+H_)R/")

	viper.SetDefault("http_service.enable", true)
	viper.SetDefault("http_service.port", 8080)
	viper.SetDefault("http_service.host", "127.0.0.1")

	viper.SetDefault("debug", false)

	viper.SetDefault("db.addr", "sqlite://.subscribe.vdb?check_same_thread=false")

	viper.SetDefault("crawler.enable", true)
	viper.SetDefault("crawler.proxies", "http://127.0.0.1:7890")
	viper.SetDefault("crawler.crawler_interval", 3600)

	viper.SetDefault("proxy_check.enable", true)
	viper.SetDefault("proxy_check.check_interval", 300)

	viper.SetDefault("proxy.enable", false)
	viper.SetDefault("proxy.mixed-port", 7890)

	viper.SetDefault("profiler.enable", false)
	viper.SetDefault("profiler.server_address", "http://localhost:4040")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("not found conf file, use default")
		} else if e, ok := err.(*os.PathError); ok {
			log.Fatalf("not find conf file in %s", e.Path)
		} else {
			log.Errorf("err:%v", err)
			return err
		}
	} else {
		// log.Infof("read config %+v", viper.AllSettings())

		// err = viper.Unmarshal(&Config)
		// if err != nil {
		// 	log.Errorf("err:%v", err)
		// 	return err
		// }

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

		log.Infof("get config %+v", Config)
	}

	if !roleList.Contains(Config.Base.Role) {
		Config.Base.Role = "Kobayashi-san"
	}

	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if Config.Profiler.Enable {
		_, _ = profiler.Start(profiler.Config{
			ApplicationName: "subscribe",
			ServerAddress:   Config.Profiler.ServerAddress,
		})
	}

	err = Ecc.Init(Config.Base.KobayashiSanHomeKey)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
