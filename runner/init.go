package runner

import (
	"errors"
	"github.com/illuminatingKong/kongming-kit/base/configx"
	"github.com/illuminatingKong/kongming-kit/base/configx/config"
	"github.com/illuminatingKong/kongming-kit/base/filesystem"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"io"
	"sync"
	"time"
)

func (o *Options) NewConfigFIle(dir, configType, name string) *Options {
	o.startConf.dir = dir
	o.startConf.configType = configType
	o.startConf.name = name
	return o
}

func (o *Options) InitBase(once *sync.Once) error {
	err := o.InitConf(once)
	if err != nil {
		return err
	}
	Logger = o.Logger
	return nil
}

func (o *Options) InitConf(once *sync.Once) error {
	var err error
	defer func() {
		if perr := recover(); perr != nil {
			o.Logger.Fatal(perr)
		}
	}()
	_config := WithConfigOption{
		ConfigDir:  o.startConf.dir,
		ConfigType: o.startConf.configType,
		ConfigName: o.startConf.name,
		ConfigIO:   o.startConf.configIO,
	}
	once.Do(func() {
		err = o.NewConf(
			WithConfX(_config))
		if err != nil {
			panic(err)
		}
	})

	o.Logger.Info("loaded config")
	return err

}

func (o *Options) WithConfWatch(once *sync.Once) {
	if o.Config == nil {
		panic(errors.New("conf is nil"))
	}

	if o.WatchConfSecond > 0 {
		ticker1 := time.NewTicker(time.Duration(o.WatchConfSecond) * time.Second)

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

func (o *Options) NewConf(opt Option) error {
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

type WithConfigOption struct {
	ConfigDir  string
	ConfigType string
	ConfigName string
	ConfigIO   io.Reader
}

func WithConfX(configOption WithConfigOption) Option {
	var newFile configx.Conf
	dir := configOption.ConfigDir
	configType := configOption.ConfigType
	name := configOption.ConfigName
	if filesystem.FileExist(dir) {
		newFile = config.NewFile(config.WithConfigsDir(config.ConfigsDir(dir)),
			config.WithConfigType(config.ConfigType(configType)),
			config.WithConfigName(config.ConfigName(name)),
			config.WithReadType(config.FileType),
		)
	} else {
		newFile = config.NewFile(config.WithConfigIO(configOption.ConfigIO),
			config.WithReadType(config.IOType))
	}

	return &option{newFile}

}
