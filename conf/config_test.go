package conf

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestInitConfig(t *testing.T) {
	err := InitConfig("../config.yaml")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
