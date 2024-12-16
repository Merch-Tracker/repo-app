package logging

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

func LogSetup(lvl string) {
	l, err := log.ParseLevel(lvl)
	if err != nil {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(
		&log.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
			},
		},
	)

	if l == log.DebugLevel {
		log.SetLevel(l)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(l)
	}

	log.SetOutput(os.Stdout)
}
