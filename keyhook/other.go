// +build !windows

package keyhook

import (
	log "github.com/sirupsen/logrus"
)

func InitKeyHook() error {
	log.Warnf("unsupported system")
	return nil
}
