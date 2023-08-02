package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type (
	Logger struct {
		*logrus.Entry
	}
	Provider interface {
		Logger() *Logger
	}
)

func NewLogger(name string, version string, level string) *Logger {
	l := logrus.New()

	switch level {
	case "error":
		l.SetLevel(logrus.ErrorLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "trace":
		l.SetLevel(logrus.TraceLevel)
	default:
		panic(fmt.Sprintf("Unknown log_level: %s used", level))
	}

	return &Logger{Entry: l.WithFields(logrus.Fields{
		"name": name, "version": version,
	})}
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Entry.Errorf(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Entry.Warnf(format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Entry.Infof(format, args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Entry.Debugf(format, args...)
}

func (l *Logger) Tracef(format string, args ...interface{}) {
	l.Entry.Tracef(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Entry.Error(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Entry.Warn(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.Entry.Info(args...)
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	ll := *l
	ll.Entry = l.Entry.WithField(key, value)
	return &ll
}

func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}

	ctx := map[string]interface{}{"message": err.Error()}

	return l.WithField("error", ctx)
}
