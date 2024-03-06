package runner

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/base/configx"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
)

var Conf configx.Conf
var Logger logx.Logger

type Option interface {
	get() interface{}
	//
}
type option struct {
	value interface{}
}

func (o *option) get() interface{} { return o.value }

type Options struct {
	OptionsCtx context.Context
	ID         string
	Name       string
	Logger     logx.Logger
	Config     configx.Conf
	startConf  startConf
	WatchConf  int
	Addr       string
}

type startConf struct {
	dir, configType, name string
}

var formatter logrusx.JsonFormatter
