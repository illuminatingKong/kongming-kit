package middleware

import (
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	ResponseError      = "error"
	ResponseData       = "response"
	WebHTTPApiError    = "web_http_api_error"
	WebHTTPApiResponse = "web_http_api_response"
	timeISO8601        = "2006-01-02T15:04:05.000Z0700"
	RequestIDName      = "requestID"
)

var sensitiveHeaders = sets.NewString("authorization", "cookie", "token", "session")
