package inter

import "github.com/illuminatingKong/kongming-kit/core/conster"

type Logger interface {
	SetApp(app AppReader) Logger
	Group(group string) Logger
	Log(severity conster.Level, message string, arguments ...interface{})
	LogWith(severity conster.Level, message string, data interface{})
}
