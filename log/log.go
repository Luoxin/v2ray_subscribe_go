package log

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Luoxin/Eutamias/utils"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

func InitLog() {
	execPath := utils.GetExecPath()
	logPath := filepath.Join(execPath, "eutamias.log")

	writer, err := rotatelogs.New(
		filepath.Join(execPath, "eutamias-%Y-%m-%d-%H.log"),
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithMaxAge(time.Hour),
		rotatelogs.WithRotationTime(time.Minute),
		// rotatelogs.WithRotationCount(24),
		// rotatelogs.WithRotationSize(100*1024*1024),
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
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
}
