package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/runner"
	"strconv"
)

const loggerKey = iota

func WithContextLogger(ctx *gin.Context) logx.Logger {
	if ctx == nil {
		return runner.Logger
	}
	l, _ := ctx.Get(strconv.Itoa(loggerKey))
	ctxLogger, ok := l.(logx.Logger)
	if ok {
		return ctxLogger
	}
	return runner.Logger
}

func NewContext(ctx *gin.Context, sub string) {
	ctx.Set(strconv.Itoa(loggerKey), WithContextLogger(ctx).Sub(sub))
}
