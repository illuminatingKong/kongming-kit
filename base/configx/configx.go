package configx

import (
	"github.com/fsnotify/fsnotify"
	"time"
)

type Conf interface {
	Load() error
	Watch() bool
	Get(str string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetIntSlice(key string) []int
	GetSizeInBytes(key string) uint
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	GetUint(key string) uint
	GetUint16(key string) uint16
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	WatchConfig()
	OnConfigChange(run func(in fsnotify.Event))
}
