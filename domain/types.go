package domain

//go:generate pie ProxyNodeList.*
type ProxyNodeList []*ProxyNode

//go:generate pie CrawlerConfList.*
type CrawlerConfList []*CrawlerConf