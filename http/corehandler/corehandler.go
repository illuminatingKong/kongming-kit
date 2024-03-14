package corehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/http/middleware"
)

type Context struct {
	Logger    logx.Logger
	Err       error
	Resp      interface{}
	RequestID string
	Page      string
	Count     string
	Limit     string
}

func NewContext(c *gin.Context) *Context {
	logger := middleware.WithContextLogger(c)
	return &Context{
		Logger:    logger,
		RequestID: c.GetString(middleware.RequestIDName),
	}
}
