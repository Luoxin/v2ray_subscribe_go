package main

import (
	"fmt"
	"time"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/eddieivan01/nic"
	log "github.com/sirupsen/logrus"

	"github.com/d-tsuji/clipboard"

	"subsrcibe/parser"
	"subsrcibe/proxycheck"
)

func main() {
	logFormatter := &nested.Formatter{
		FieldsOrder: []string{
			log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
			log.FieldKeyFunc, log.FieldKeyMsg,
		},
		TimestampFormat:  time.RFC3339,
		HideKeys:         true,
		NoFieldsSpace:    true,
		NoUppercaseLevel: true,
		TrimMessages:     true,
		CallerFirst:      true,
	}

	log.SetFormatter(logFormatter)
	log.SetReportCaller(true)

	err := func() error {

		return nil
	}()
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}
}
