package corehandler

import (
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/http/middleware"
)

func WebHttpApiResponse(c *gin.Context, ctx *Context) {
	c.Header("core-requestID", ctx.RequestID)
	if ctx.Err != nil {
		ctx.Logger.Errorf("requestID: %s  Error: %+v", ctx.RequestID, ctx.Err)
		c.Set(middleware.WebHTTPApiError, ctx.Err)
		c.Abort()
		return
	}

	if ctx.Resp != nil {
		ctx.Logger.Infof("requestID: %s  Response: %+v", ctx.RequestID, ctx.Resp)
		realResp := responseHelper(ctx.Resp)
		c.Set(middleware.WebHTTPApiResponse, realResp)
	}
}

//func responseHelper(response interface{}) interface{} {
//	switch response.(type) {
//	case string, []byte:
//		return response
//	}
//	return response
//}
