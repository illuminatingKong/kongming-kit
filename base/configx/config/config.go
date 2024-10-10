package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/illuminatingKong/kongming-kit/base/configx"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
	"github.com/spf13/viper"
	"io"
)

type ConfigsDir string
type ConfigName string
type ConfigType string
type ReadType string
type ConfigIO io.Reader

var (
	FileType    ReadType = "file"
	IOType      ReadType = "io"
	FindNotType          = errors.New("find not config type")
)

type Config struct {
	ReadType   ReadType
	ConfigIO   ConfigIO
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

func WithReadType(ReadType ReadType) Option {
	return &option{ReadType}
}

func WithConfigIO(ConfigIO io.Reader) Option {
	return &option{ConfigIO}
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
	case ReadType:
		c.ReadType = val
	case ConfigIO:
		c.ConfigIO = val
	}
}

func (c *Config) Load() error {
	// Default
	if c.ReadType == "" {
		return FindNotType
	}
	if c.ReadType == FileType || len(string(c.ConfigsDir)) > 0 {
		c.Viper.SetConfigName(string(c.ConfigName))
		c.Viper.SetConfigType(string(c.Type))
		c.SetConfigxPath(string(c.ConfigsDir))
		return c.Viper.ReadInConfig()
	}
	c.Viper.SetConfigType(string(c.Type))
	return c.Viper.ReadConfig(c.ConfigIO)
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
