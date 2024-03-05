package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/illuminatingKong/kongming-kit/base/configx"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
	"github.com/spf13/viper"
)

type ConfigsDir string
type ConfigName string
type ConfigType string

type Config struct {
	ConfigsDir ConfigsDir
	ConfigName ConfigName
	Type       ConfigType
	Logger     logx.Logger
	*viper.Viper
}

var formatter logrusx.JsonFormatter

func NewFile(options ...Option) configx.Conf {
	var c = &Config{
		Logger: logrusx.New(logrusx.WithFormatter(formatter)),
		Viper:  viper.New(),
	}

	for _, opt := range options {
		processOptions(c, opt.get())
	}

	return c
}

func WithConfigsDir(ConfigsDir ConfigsDir) Option {
	return &option{ConfigsDir}
}

func WithConfigName(ConfigName ConfigName) Option {
	return &option{ConfigName}
}

func WithConfigType(ConfigType ConfigType) Option {
	return &option{ConfigType}
}

func processOptions(c *Config, opt interface{}) {
	switch val := opt.(type) {
	case ConfigsDir:
		c.ConfigsDir = val
	case ConfigName:
		c.ConfigName = val
	case ConfigType:
		c.Type = val
	}
}

func (c *Config) Load() error {
	c.Viper.SetConfigName(string(c.ConfigName))
	c.Viper.SetConfigType(string(c.Type))
	c.SetConfigxPath(string(c.ConfigsDir))
	if err := c.Viper.ReadInConfig(); err != nil { // viper解析配置文件
		panic(err)
	}
	return nil
}

func (c *Config) Watch() bool {
	var v bool
	c.Viper.WatchConfig()
	c.Viper.OnConfigChange(func(e fsnotify.Event) {
		c.Logger.Infof("Config file changed: %+v", e.Name)
		v = true
	})
	return v
}
