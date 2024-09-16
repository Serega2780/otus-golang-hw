package logger

import (
	"io"
	"os"

	cfg "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func New(conf *cfg.LoggerConf) *Logger {
	log := logrus.New()
	l, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.SetLevel(l)
	switch conf.Format {
	case "text":
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.Errorf("Unknown logging formatter received %s", conf.Format)
		os.Exit(1)
	}

	if conf.LogToFile {
		file, err := os.OpenFile(conf.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o666)
		if err != nil {
			log.Errorf("Unable to create log file %s", err.Error())
			os.Exit(1)
		}

		switch {
		case conf.LogToConsole && conf.LogToFile:
			log.SetOutput(io.MultiWriter(file, os.Stdout))
		case conf.LogToConsole:
			log.SetOutput(io.MultiWriter(os.Stdout))
		case conf.LogToFile:
			log.SetOutput(io.MultiWriter(file))
		}
	}

	return &Logger{log: log}
}

func (l *Logger) GetWriter() io.Writer {
	return l.log.Out
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log.Warn(args...)
}
