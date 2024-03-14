package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/http/errorsx"
	"github.com/illuminatingKong/kongming-kit/http/webapi"
	"net/http"
)

// Response handle response
func Response() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer handleResponse(c)
		c.Next()
	}
}

func handleResponse(c *gin.Context) {

	if v, ok := c.Get(ResponseError); ok {
		c.JSON(errorsx.ErrorMessage(v.(error)))
		return
	}

	if v, ok := c.Get(ResponseData); ok {
		setResponse(v, c)
	}

	if v, ok := c.Get(WebHTTPApiError); ok {
		c.JSON(errorsx.ErrorMessage(v.(error)))
		return
	}
	if v, ok := c.Get(WebHTTPApiResponse); ok {
		c.JSON(webapi.WebHttpApiResponseHandler(v))
	}

	// skip if response or status is already set
	if c.Writer.Written() || c.Writer.Status() != http.StatusOK {
		return
	}

}

func setResponse(resp interface{}, c *gin.Context) {
	switch r := resp.(type) {
	case string:
		c.String(200, r)
	case []byte:
		c.String(200, string(r))
	default:
		c.JSON(200, format(r))
	}
}

func format(resp interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": 200,
		"data": resp,
	}
}
