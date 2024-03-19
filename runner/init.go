package runner

import (
	"context"
	"errors"
	"github.com/illuminatingKong/kongming-kit/base/configx"
	"github.com/illuminatingKong/kongming-kit/base/configx/config"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"sync"
	"time"
)

func (o *Options) NewConfig(dir, configType, name string) *Options {
	o.startConf.dir = dir
	o.startConf.configType = configType
	o.startConf.name = name
	return o
}

func (o *Options) InitBase(ctx context.Context, once *sync.Once) error {
	err := o.InitConf(ctx, once)
	if err != nil {
		return err
	}
	Logger = o.Logger
	return nil
}

func (o *Options) InitConf(ctx context.Context, once *sync.Once) error {
	var err error
	defer func() {
		if perr := recover(); perr != nil {
			o.Logger.Fatal(perr)
		}
	}()

	once.Do(func() {
		err = o.NewConf(ctx,
			WithConfX(o.startConf.dir, o.startConf.configType, o.startConf.name))
		if err != nil {
			panic(err)
		}
	})

	o.Logger.Info("loaded config")
	return err

}

func (o *Options) WithConfWatch(ctx context.Context, once *sync.Once) {
	if o.Config == nil {
		panic(errors.New("conf is nil"))
	}

	if o.WatchConf > 0 {
		ticker1 := time.NewTicker(time.Duration(o.WatchConf) * time.Second)

		go func(t *time.Ticker, o *Options) {
			for {
				<-t.C
				once.Do(func() {
					change := o.Config.Watch()
					if change {
						Conf = o.Config
					}
				})

			}
		}(ticker1, o)
	}

}

func (o *Options) NewConf(ctx context.Context, opt Option) error {
	var err error
	o.Config = opt.get().(configx.Conf)
	err = o.Config.Load()
	if err != nil {
		return err
	}
	Conf = o.Config
	return err
}

func GetConf() configx.Conf {
	return Conf
}

func GetLogger() logx.Logger {
	if Logger == nil {
		panic(errors.New("logger is nil"))
	}
	return Logger
}

func WithConfX(dir, configType, name string) Option {
	newFile := config.NewFile(config.WithConfigsDir(config.ConfigsDir(dir)),
		config.WithConfigType(config.ConfigType(configType)),
		config.WithConfigName(config.ConfigName(name)))
	return &option{newFile}

}
