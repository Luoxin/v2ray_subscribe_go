package main

import "subsrcibe/domain"

// //go:generate pie CrawlerConfList.*
type CrawlerConfList []*domain.CrawlerConf

// //go:generate pie ProxyNodeList.*
type ProxyNodeList []*domain.ProxyNode
