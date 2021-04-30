// +build darwin

package notify

import (
	"fmt"

	"github.com/deckarep/gosx-notifier"
	log "github.com/sirupsen/logrus"
)

func Msg(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	log.Infof("notify:%v", text)

	note := gosxnotifier.NewNotification(text)
	note.Title = "eutamias"
	note.Subtitle = "eutamias"

	_ = note.Push()
}
