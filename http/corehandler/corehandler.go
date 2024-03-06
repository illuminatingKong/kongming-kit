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

func JSONResponse(c *gin.Context, ctx *Context) {
	c.Header("NioRequestID", ctx.RequestID)
	if ctx.Err != nil {
		ctx.Logger.Errorf("RequestID: %s  Error: %+v", ctx.RequestID, ctx.Err)
		c.Set(middleware.ResponseError, ctx.Err)
		c.Abort()
		return
	}

	if ctx.Resp != nil {
		ctx.Logger.Infof("RequestID: %s  Response: %+v", ctx.RequestID, ctx.Resp)
		realResp := responseHelper(ctx.Resp)
		c.Set(middleware.ResponseData, realResp)
	}
}

func responseHelper(response interface{}) interface{} {
	switch response.(type) {
	case string, []byte:
		return response
	}
	return response
}
