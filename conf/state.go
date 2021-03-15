package conf

import (
	"github.com/bluele/gcache"

	"subscribe/utils"
)

var Cache = gcache.New(20).LRU().Build()
var Ecc = utils.NewEcc()
