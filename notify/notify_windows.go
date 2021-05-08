// +build windows

package notify

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/toast.v1"
)

func Msg(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)

	n := &toast.Notification{
		AppID:    "eutamias",
		Title:    "eutamias",
		Message:  text,
		Duration: toast.Short,
	}

	log.Infof("notify:%v", text)

	err := n.Push()
	if err != nil {
		log.Errorf("err:%v", err)
	}
}
