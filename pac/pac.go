package pac

import (
	"sync"
)

type pac struct {
}

var Pac *pac
var lock sync.RWMutex

func init() {
	Pac = NewPac()
}

func NewPac() *pac {
	return &pac{}
}

//func (p *pac) UpdatePac()  {
//	lock.Lock()
//	lock.Unlock()
//
//	nic.Get("http://cdn.jsdelivr.net/gh/gfwlist/gfwlist@master/gfwlist.txt")
//}
