package eutamias

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Luoxin/Eutamias/geolite"
	"github.com/Luoxin/Eutamias/utils"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/db"
	"github.com/Luoxin/Eutamias/task"
	"github.com/Luoxin/Eutamias/webservice"
)

func init() {
	InitLog()
}

func Init(configFilePatch string) error {
	err := conf.InitConfig(configFilePatch)
	if err != nil {
		log.Fatalf("init config err:%v", err)
		return err
	}

	log.Info("init conf success")

	err = geolite.InitGeoLite()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info("init geolite success")

	err = db.InitDb(conf.Config.Db.Addr)
	if err != nil {
		log.Fatalf("init db err:%v", err)
		return err
	}

	log.Info("init db success")

	err = task.InitWorker()
	if err != nil {
		log.Fatalf("init work err:%v", err)
		return err
	}

	log.Info("init worker success")

	err = webservice.InitHttpService()
	if err != nil {
		log.Fatalf("init http service err:%v", err)
		return err
	}

	log.Info("init http service success")

	return nil
}

func Start() {
	c := make(chan bool)

	err := Init("")
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}

	<-c
}

func InitLog() {
	execPath := utils.GetExecPath()
	logPath := filepath.Join(execPath, "eutamias.log")

	writer, err := rotatelogs.New(
		filepath.Join(execPath, "eutamias-%Y-%m-%d-%H.log"),
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Hour),
		rotatelogs.WithRotationTime(time.Minute),
	)
	if err != nil {
		log.Fatalf("err:%v", err)
	}

	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.InfoLevel:  writer,
			log.WarnLevel:  writer,
			log.ErrorLevel: writer,
			log.FatalLevel: writer,
			log.PanicLevel: writer,
		},
		&nested.Formatter{
			FieldsOrder: []string{
				log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
				log.FieldKeyFunc, log.FieldKeyMsg,
			},
			CustomCallerFormatter: func(f *runtime.Frame) string {
				return fmt.Sprintf("(%s %s:%d)", f.Function, path.Base(f.File), f.Line)
			},
			TimestampFormat:  time.RFC3339,
			HideKeys:         true,
			NoFieldsSpace:    true,
			NoUppercaseLevel: true,
			TrimMessages:     true,
			CallerFirst:      true,
		},
	))
	// log.AddHook(lfshook.NewHook(
	// 	lfshook.WriterMap{
	// 		log.TraceLevel: os.Stdout,
	// 		log.DebugLevel: os.Stdout,
	// 		log.InfoLevel:  os.Stdout,
	// 		log.WarnLevel:  os.Stdout,
	// 		log.ErrorLevel: os.Stdout,
	// 		log.FatalLevel: os.Stdout,
	// 		log.PanicLevel: os.Stdout,
	// 	},
	// 	&nested.Formatter{
	// 		FieldsOrder: []string{
	// 			log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
	// 			log.FieldKeyFunc, log.FieldKeyMsg,
	// 		},
	// 		CustomCallerFormatter: func(f *runtime.Frame) string {
	// 			return fmt.Sprintf("(%s %s:%d)", f.Function, path.Base(f.File), f.Line)
	// 		},
	// 		TimestampFormat:  time.RFC3339,
	// 		HideKeys:         true,
	// 		NoFieldsSpace:    true,
	// 		NoUppercaseLevel: true,
	// 		TrimMessages:     true,
	// 		CallerFirst:      true,
	// 	},
	// ))

	log.SetFormatter(&nested.Formatter{
		FieldsOrder: []string{
			log.FieldKeyTime, log.FieldKeyLevel, log.FieldKeyFile,
			log.FieldKeyFunc, log.FieldKeyMsg,
		},
		CustomCallerFormatter: func(f *runtime.Frame) string {
			return fmt.Sprintf("(%s %s:%d)", f.Function, path.Base(f.File), f.Line)
		},
		TimestampFormat:  time.RFC3339,
		HideKeys:         true,
		NoFieldsSpace:    true,
		NoUppercaseLevel: true,
		TrimMessages:     true,
		CallerFirst:      true,
	})
	log.SetReportCaller(true)
}
