package conf

import (
	"github.com/bluele/gcache"

	"github.com/luoxin/subscribe/utils"
)

var Cache = gcache.New(20).LRU().Build()
var Ecc = utils.NewEcc()
