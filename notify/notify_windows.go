// +build windows

package notify

import (
	"fmt"

	eutamias "github.com/Luoxin/Eutamias"
	log "github.com/sirupsen/logrus"
	"gopkg.in/toast.v1"
)

func Msg(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)

	n := &toast.Notification{
		AppID:    eutamias.ServiceName,
		Title:    eutamias.ServiceName,
		Message:  text,
		Duration: toast.Short,
	}

	log.Infof("notify:%v", text)

	_ = n.Push()
}
