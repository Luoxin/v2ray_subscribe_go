package main

type CrawlerConf struct {
	ID        uint64 `json:"id" gorm:"primary_key"`
	CreatedAt uint32 `json:"created_at"`
	UpdatedAt uint32 `json:"updated_at"`
	CrawlUrl  string `json:"crawl_url" gorm:"type:varchar(1000);unique_index:idx_crawl_url"`
	CrawlType uint32 `json:"crawl_type"`
	Rule      string `json:"rule" gorm:"type:json"`
	IsClosed  bool   `json:"is_closed"`
	NextAt    uint32 `json:"next_at" gorm:"index:idx_next_crawl_at"`
	Interval  uint32 `json:"interval"`
	Note      string `json:"note" gorm:"type:varchar(1000)"`
}
