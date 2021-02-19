package conf

import (
	"github.com/bluele/gcache"
)

var Cache = gcache.New(20).LRU().Build()
