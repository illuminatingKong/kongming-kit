package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/uuid"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

func RequestID(logger logx.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		reqID := uuid.ID()
		NewContext(c, reqID)
		c.Set(RequestIDName, reqID)

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		var body []byte
		headers := make(map[string]string)
		// request body is a ReadCloser, it can be read only once.
		if c.Request != nil {
			if c.Request.Body != nil {
				var buf bytes.Buffer
				tee := io.TeeReader(c.Request.Body, &buf)
				body, _ = ioutil.ReadAll(tee)
				c.Request.Body = ioutil.NopCloser(&buf)
			}

			for k := range c.Request.Header {
				if sensitiveHeaders.Has(strings.ToLower(k)) {
					continue
				}
				headers[k] = c.GetHeader(k)
			}
		}

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		logger.Infof("status: %d "+
			"method: %s "+
			"path: %s "+
			"body: %s "+
			"headers: %v "+
			"clientIP: %s "+
			"user-agent: %s "+
			"requestID: %s "+
			"start: %s "+
			"latency: %s "+
			"error: %s", c.Writer.Status(), c.Request.Method, path, string(body),
			headers, c.ClientIP(), c.Request.UserAgent(), c.GetString(RequestIDName),
			start.Format(timeISO8601), latency.String(),
			c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}

}
