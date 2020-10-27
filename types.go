package main

import "subsrcibe/subscription"

// //go:generate pie CrawlerConfList.*
type CrawlerConfList []*subscription.CrawlerConf

// //go:generate pie ProxyNodeList.*
type ProxyNodeList []*subscription.ProxyNode
