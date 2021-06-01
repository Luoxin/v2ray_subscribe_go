package conf

import (
	"github.com/bluele/gcache"

	"github.com/Luoxin/Eutamias/utils"
)

var Cache = gcache.New(1024).LRU().Build()
var Ecc = utils.NewEcc()
