package log

import (
	"github.com/illuminatingKong/kongming-kit/core/conster"
	"github.com/illuminatingKong/kongming-kit/core/container"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	name string
	app  container.AppReader
	*logrus.Entry
}

func (l Logger) SetApp(app container.AppReader) container.Logger {
	l.app = app
	return l
}

func (l Logger) Group(group string) container.Logger {
	l.Entry.WithField("group", group)
	return l
}

func (l Logger) Log(severity conster.Level, message string, arguments ...interface{}) {
	l.Entry.Level = logrus.Level(severity)
	l.Entry.WithField("message", message)
	l.Entry.Info(arguments...)
}

func (l Logger) LogWith(severity conster.Level, message string, data interface{}) {
	l.Entry.Level = logrus.Level(severity)
	l.Entry.WithField(message, data)
}
