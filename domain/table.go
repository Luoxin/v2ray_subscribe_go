package domain

type CrawlType uint32

const (
	CrawlTypeNil           CrawlType = 0
	CrawlTypeSubscription  CrawlType = 1
	CrawlTypeXpath         CrawlType = 2
	CrawlTypeFuzzyMatching CrawlType = 3
	CrawlTypeClashProxies  CrawlType = 4
)

type CrawlerUse uint32

const (
	CrawlerUseNil     CrawlerUse = 0
	CrawlerUseGFW     CrawlerUse = 1
	CrawlerUseNetEase CrawlerUse = 2 // 网易云
)

type CrawlerConf struct {
	Id             uint64     `json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt      uint32     `json:"created_at,omitempty" gorm:"autoUpdateTime:autoCreateTime"`
	UpdatedAt      uint32     `json:"updated_at,omitempty" gorm:"autoUpdateTime:autoUpdateTime"`
	CrawlerFeature string     `json:"crawler_features,omitempty" gorm:"type:varchar(700);unique_index:idx_crawl_url;comment:地址的特征，因为mysql索引不支持太长的，所以用sha_512_做唯一索引"`
	CrawlUrl       string     `json:"crawl_url,omitempty" gorm:"type:text;comment:抓取的地址"`
	CrawlType      CrawlType  `json:"crawl_type,omitempty"`
	CrawlerUse     CrawlerUse `json:"crawler_use,omitempty"`
	// @gorm: type:json
	Rule *CrawlerConf_Rule `json:"rule,omitempty" gorm:"type:json"`
	// @grom: index: idx_next_crawl_at
	IsClosed bool `json:"is_closed,omitempty" gorm:"index:idx_next_crawl_at"`
	// @grom: index: idx_next_crawl_at
	NextAt   uint32 `json:"next_at,omitempty" gorm:"index:idx_next_crawl_at"`
	Interval uint32 `json:"interval,omitempty"`
	// @v: max=100
	Note string `json:"note,omitempty"`
}

type CrawlerConf_Rule struct {
	UseProxy bool   `json:"use_proxy,omitempty"`
	Xpath    string `json:"xpath,omitempty"`
}

type ProxyNodeType uint32

const (
	ProxyNodeType_ProxyNodeTypeNil    ProxyNodeType = 0
	ProxyNodeType_ProxyNodeTypeVmess  ProxyNodeType = 1
	ProxyNodeType_ProxyNodeTypeTrojan ProxyNodeType = 2
	ProxyNodeType_ProxyNodeTypeVless  ProxyNodeType = 3
	ProxyNodeType_ProxyNodeTypeSS     ProxyNodeType = 4
	ProxyNodeType_ProxyNodeTypeSSR    ProxyNodeType = 5
	ProxyNodeType_ProxyNodeTypeSocket ProxyNodeType = 6
	ProxyNodeType_ProxyNodeTypeHttp   ProxyNodeType = 7
	ProxyNodeType_ProxyNodeTypeSnell  ProxyNodeType = 8
)

type ProxyUse uint32

const (
	ProxyUseNil     ProxyUse = 0
	ProxyUseGFW     ProxyUse = 1
	ProxyUseNetEase ProxyUse = 2 // 网易云
)

type ProxyNode struct {
	Id        uint64 `json:"id,omitempty"`
	CreatedAt uint32 `json:"created_at,omitempty" gorm:"autoUpdateTime:autoCreateTime"`
	UpdatedAt uint32 `json:"updated_at,omitempty" gorm:"autoUpdateTime:autoUpdateTime"`

	UrlFeature        string        `json:"url_feature,omitempty" gorm:"type:varchar(700);not_null;unique_index:idx_node_url;comment:节点的地址，因为mysql索引不支持太长的，所以用sha_512_做唯一索引"`
	Url               string        `json:"url,omitempty" gorm:"type:text;comment:节点的地址"`
	ProxyNodeType     ProxyNodeType `json:"proxy_node_type,omitempty"`
	ProxySpeed        float64       `json:"proxy_speed,omitempty"`
	ProxyNetworkDelay float64       `json:"proxy_network_delay,omitempty"`
	NextCheckAt       uint32        `json:"next_check_at,omitempty" gorm:"index:idx_next_check_at"`
	CheckInterval     uint32        `json:"check_interval,omitempty"`
	CrawlId           uint64        `json:"crawl_id,omitempty"`
	DeathCount        uint32        `json:"death_count,omitempty" gorm:"index:idx_alive"`
	IsClose           bool          `json:"is_close,omitempty" gorm:"index:idx_next_check_at;index:idx_alive"`
	LastCrawlerAt     uint32        `json:"last_crawler_at,omitempty"`
	AvailableCount    uint64        `json:"available_count,omitempty" gorm:"index:idx_alive"`
	ProxyUse          ProxyUse      `json:"proxy_use,omitempty"`
}

type TohruFeed struct {
	Id        uint64 `json:"id,omitempty"`
	CreatedAt uint32 `json:"created_at,omitempty" gorm:"autoUpdateTime:autoCreateTime"`
	UpdatedAt uint32 `json:"updated_at,omitempty" gorm:"autoUpdateTime:autoUpdateTime"`

	UserId       string `json:"user_id,omitempty" gorm:"unique_index:idx_user_id;not null;comment:用户的唯一标识"`
	UserPassword string `json:"user_password,omitempty" gorm:"unique_index:idx_user_id;not_null"`

	UpCount   uint64 `json:"up_count,omitempty" gorm:"comment:用户上传的数量"`
	LastIp    string `json:"last_ip,omitempty" gorm:"type:varchar(100)"`
	IsDisable string `json:"is_disable,omitempty" gorm:"comment:是否禁用"`
}
