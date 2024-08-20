package log

import (
	"github.com/illuminatingKong/kongming-kit/base/msg"
	"github.com/sirupsen/logrus"
)

func (o *option) get() interface{} { return o.value }

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

func (l Logger) AddHook(hook logrus.Hook) {
	l.Entry.Logger.AddHook(hook)
}

func (l Logger) WithFieldsX(f func(args ...interface{}), m string, keyvals ...interface{}) {
	f(msg.LogMessage(m, keyvals...))
}
