package runner

import (
	"context"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
)

func NewContainer(name, instance string) *Options {
	o := &Options{
		OptionsCtx:      context.Background(),
		Name:            name,
		Logger:          logrusx.New(logrusx.WithFormatter(formatter)),
		Config:          nil,
		WatchConfSecond: 2,
		ID:              instance,
	}
	return o
}
