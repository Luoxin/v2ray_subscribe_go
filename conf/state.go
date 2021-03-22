package conf

import (
	"github.com/bluele/gcache"

	"github.com/luoxin/v2ray_subscribe_go/utils"
)

var Cache = gcache.New(20).LRU().Build()
var Ecc = utils.NewEcc()
