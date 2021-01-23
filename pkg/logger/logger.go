package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
}

var instance Logger

func init() {
	InitLogger("debug")
}

func Debug(args ...interface{}) {
	instance.Debug(args)
}

func Debugf(template string, args ...interface{}) {
	instance.Debugf(template, args...)
}

func Error(args ...interface{}) {
	instance.Error(args)
}

func Errorf(template string, args ...interface{}) {
	instance.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	instance.Fatal(args)
}

func Fatalf(template string, args ...interface{}) {
	instance.Fatalf(template, args...)
}

func Info(args ...interface{}) {
	instance.Info(args)
}

func Infof(template string, args ...interface{}) {
	instance.Infof(template, args...)
}

func Warn(args ...interface{}) {
	instance.Warn(args)
}

func Warnf(template string, args ...interface{}) {
	instance.Warnf(template, args...)
}

func InitLogger(logLevel string) {
	initLogrus(logLevel)
}

func initLogrus(logLevel string) {
	logEntry := logrus.New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Error("Error parsing log level, using: info")
		level = logrus.InfoLevel
	}
	logEntry.Level = level
	Formatter := new(logrus.TextFormatter)
	Formatter.TimestampFormat = "20060102T15:04:05.999"
	Formatter.FullTimestamp = true
	if level == logrus.DebugLevel {
		Formatter.DisableQuote = true
	}
	logEntry.SetFormatter(Formatter)
	instance = logEntry
}

// Get original instance.
func Logrus() *logrus.Entry {
	if entry, ok := instance.(*logrus.Entry); ok {
		return entry
	}
	return nil
}
