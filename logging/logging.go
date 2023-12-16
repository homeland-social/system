package logging

import (
	"fmt"

	"log/syslog"
)

var logger *syslog.Writer = nil

func getLogger() (*syslog.Writer, error) {
	var err error

	if logger == nil {
		logger, err = syslog.New(syslog.LOG_DEBUG, "system api")
	}
	return logger, err
}

func Debug(m string, v ...any) {
	log, err := getLogger()
	if err != nil {
		return
	}
	s := fmt.Sprintf(m, v...)
	log.Debug(s)
}

func Error(m string, v ...any) {
	log, err := getLogger()
	if err != nil {
		return
	}
	s := fmt.Sprintf(m, v...)
	log.Alert(s)
}
