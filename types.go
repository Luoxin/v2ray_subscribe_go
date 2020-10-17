package main

import "subsrcibe/subscribe"

//go:generate pie CrawlerConfList.*
type CrawlerConfList []*subscribe.CrawlerConf
