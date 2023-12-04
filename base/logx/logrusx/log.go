package logrusx

import (
	"errors"
	"fmt"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/msg"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	name string
	*logrus.Entry
}

type Option interface {
	get() interface{}
}

type setNameOption string

type JsonFormatter logrus.Formatter

type option struct{ value interface{} }

func (o *option) get() interface{} { return o.value }

type PathMap map[logrus.Level]string

func WithFormatter(l JsonFormatter) Option {
	l = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "@level",
		},
	}
	return &option{l}
}

func (l Logger) AddHook(hook interface{}) error {
	var err error
	defer func() {
		if p := recover(); p != nil {
			err = errors.New(fmt.Sprintf("Add hook has error: %+v", p))
		}
	}()
	h := hook.(logrus.Hook)
	logrus.AddHook(h)
	return err
}

func (l Logger) WithFieldsX(f func(args ...interface{}), m string, keyvals ...interface{}) {
	f(msg.LogMessage(m, keyvals...))
}

func New(options ...Option) logx.Logger {
	log := logrus.New()
	logger := &Logger{"", logrus.NewEntry(log)}
	for _, opt := range options {
		processOptions(log, logger, opt.get())
	}
	return logger
}

func (l Logger) Sub(name string) logx.Logger {
	if len(l.name) > 0 {
		name = l.name + "." + name
	}
	return &Logger{name, l.Entry.WithField("module", name)}
}

func processOptions(logr *logrus.Logger, logger *Logger, opt interface{}) {
	switch val := opt.(type) {
	case setNameOption:
		logger.name = string(val)
	case logrus.Level:
		logr.SetLevel(val)
	case JsonFormatter:
		logr.SetFormatter(val)
	}
}
