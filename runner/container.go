package runner

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
)

func NewContainer(name, addr string) *Options {
	o := &Options{
		OptionsCtx: context.Background(),
		Name:       name,
		Logger:     logrusx.New(logrusx.WithFormatter(formatter)),
		Config:     nil,
		WatchConf:  2,
		Addr:       addr,
	}
	return o
}

var CLog = func() logx.Logger {
	return logrusx.New(logrusx.WithFormatter(formatter))
}
