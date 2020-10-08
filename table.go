package main

type CrawlerConf struct {
	ID        uint64 `json:"id"`
	CreatedAt uint32 `json:"created_at"`
	UpdatedAt uint32 `json:"updated_at"`
	CrawlURL  string `json:"crawl_url" gorm:"type:varchar(1000)"`
	CrawlType uint32 `json:"crawl_type"`
	Rule      string `json:"rule" gorm:"type:json"`
	IsClosed  bool   `json:"is_closed"`
	NextAt    uint32 `json:"next_at"`
	Interval  uint32 `json:"interval"`
	Note      string `json:"note" gorm:"type:varchar(1000)"`
}
