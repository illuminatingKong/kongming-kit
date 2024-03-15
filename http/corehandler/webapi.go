package corehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/http/middleware"
)

func WebHttpApiResponse(c *gin.Context, ctx *Context) {
	requestID := ctx.RequestID
	requestFullPath := c.FullPath()
	c.Header("core-requestID", requestID)
	if ctx.Err != nil {
		ctx.Logger.Errorf("requestID: %s, error: %+v, requestURL: %s", requestID, ctx.Err, requestFullPath)
		c.Set(middleware.WebHTTPApiError, ctx.Err)
		c.Abort()
		return
	}

	if ctx.Resp != nil {
		ctx.Logger.Infof("requestID: %s, response: %+v, requestURL: %s", requestID, ctx.Resp, requestFullPath)
		realResp := responseHelper(ctx.Resp)
		c.Set(middleware.WebHTTPApiResponse, realResp)
	}
}
