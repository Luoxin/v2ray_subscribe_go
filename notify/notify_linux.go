// +build linux

package notify

import (
	log "github.com/sirupsen/logrus"
)

func Msg(text string) {
	log.Infof("notify:%v", text)
}
