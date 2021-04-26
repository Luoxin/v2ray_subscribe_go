// +build !windows

package notify

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func Msg(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	log.Infof("notify:%v", text)
}
