package log

import (
	"github.com/sirupsen/logrus"
)

type Option interface {
	get() interface{}
}

type setAppName string

type JsonFormatter logrus.Formatter

type option struct{ value interface{} }
