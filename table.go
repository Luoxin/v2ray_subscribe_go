package main

type CrawlType uint32

const (
	CrawlTypeNil = iota
	CrawlTypeSubscription
	CrawlTypeXpath
)

type CrawlerConf_Rule struct {
	UseProxy bool   `json:"use_proxy"`
	Xpath    string `json:"xpath"`
}

type CrawlerConf struct {
	Id        uint64           `json:"id" gorm:"primary_key"`
	CreatedAt uint32           `json:"created_at"`
	UpdatedAt uint32           `json:"updated_at"`
	CrawlUrl  string           `json:"crawl_url" gorm:"type:varchar(1000);unique_index:idx_crawl_url"`
	CrawlType CrawlType        `json:"crawl_type"`
	Rule      CrawlerConf_Rule `json:"rule" gorm:"type:json"`
	IsClosed  bool             `json:"is_closed"`
	NextAt    uint32           `json:"next_at" gorm:"index:idx_next_crawl_at"`
	Interval  uint32           `json:"interval"`
	Note      string           `json:"note" gorm:"type:varchar(1000)"`
}

type ProxyNodeType uint32

const (
	ProxyNodeTypeNil = iota
	ProxyNodeTypeVmess
)

type ProxyNode struct {
	Id                uint64        `json:"id" gorm:"primary_key"`
	CreatedAt         uint32        `json:"created_at"`
	UpdatedAt         uint32        `json:"updated_at"`
	Url               string        `json:"url" gorm:"type:varchar(1000)"`
	NodeType          ProxyNodeType `json:"node_type"`
	NodeDetail        string        `json:"node_detail"`
	ProxySpeed        string        `json:"proxy_speed"`
	ProxyNetworkDelay uint32        `json:"proxy_network_delay"`
	NextCheckAt       uint32        `json:"next_check_at" gorm:"index:idx_next_check_at"`
	CheckInterval     uint32        `json:"check_interval"`
	CrawlId           uint64        `json:"crawl_id"`
	IsClosed          string        `json:"is_closed" gorm:"index:idx_next_check_at_next_at"`
}

type VmessNode struct {
	V    string `json:"v"`
	Id   string `json:"id"`
	PS   string `json:"ps"`
	Add  string `json:"add"`
	Aid  string `json:"aid"`
	Net  string `json:"net"`
	TLS  string `json:"tls"`
	Host string `json:"host"`
	Path string `json:"path"`
	Port string `json:"port"`
	Type string `json:"type"`
}
