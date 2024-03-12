package rollout

import (
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
)

var logger logx.Logger

func InitLog() logx.Logger {
	logger = logrusx.New()
	return logger
}

func Logger() logx.Logger {
	return getLogger()
}

func getLogger() logx.Logger {
	if logger == nil {
		panic("logger is not initialized yet!")
	}

	return logger
}
