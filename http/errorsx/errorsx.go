package errorsx

import (
	"fmt"
	"regexp"
)

var (
	ErrInternalError = NewHTTPError(500, "Internal Error")
	ErrInvalidParam  = NewHTTPError(400, "Bad Request")
	// ErrUnauthorized ...
	ErrUnauthorized = NewHTTPError(401, "Unauthorized")
	// ErrForbidden ...
	ErrForbidden = NewHTTPError(403, "Forbidden")
	// ErrNotFound ...
	ErrNotFound          = NewHTTPError(404, "Request Not Found")
	InvalidDataInRequest = "invalid data json in request"
	// ResponseError  export common response error, can use addErr to add error replace default error
	RespnseError = NewHTTPError(4000, "Response Error")
)

type IHTTPError interface {
	Code() int
	Error() string
	Message() string
	Desc() string
	Extra() map[string]interface{}
	Data() interface{}
}

// HTTPError ...
type HTTPError struct {
	code  int
	err   string
	desc  string
	extra map[string]interface{}
	data  interface{}
}

// NewHTTPError ...
func NewHTTPError(code int, errStr string, args ...string) *HTTPError {

	var desc string
	if len(args) > 0 {
		desc = args[0]
	}

	return &HTTPError{
		code: code,
		err:  errStr,
		desc: desc,
	}
}

// Code ...
func (e *HTTPError) Code() int {
	return e.code
}

// Error ...
func (e *HTTPError) Error() string {
	return fmt.Sprintf("%s: %s", e.err, e.desc)
}

func (e *HTTPError) Message() string {
	return e.err
}

// Desc ...
func (e *HTTPError) Desc() string {
	return e.desc
}

// Extra ...
func (e *HTTPError) Extra() map[string]interface{} {
	extra := map[string]interface{}{}

	for k, v := range e.extra {
		extra[k] = v
	}

	return extra
}

func (e *HTTPError) Data() interface{} {
	return e.data
}

// AddDesc ...
func (e *HTTPError) AddDesc(desc string) *HTTPError {
	// set default description error
	e.desc = desc

	if matched, _ := regexp.MatchString(".*E11000 duplicate.*", desc); matched {
		e.desc = "mongo duplicate key error"
	}

	return e
}

// AddErr ...
func (e *HTTPError) AddErr(err error) *HTTPError {
	e.err = err.Error()
	return e
}

func (e *HTTPError) SetData(data interface{}) *HTTPError {
	e.data = data
	return e
}

// NewWithDesc ...
func NewWithDesc(e error, desc string) error {
	if v, ok := e.(*HTTPError); ok {
		err := *v
		err.desc = desc
		return &err
	}

	return e
}

// NewWithExtras ...
func NewWithExtras(e error, desc string, extra map[string]interface{}) error {
	if v, ok := e.(*HTTPError); ok {
		err := *v
		err.desc = desc
		err.extra = extra
		return &err
	}

	return e
}

// ErrorMessage returns the code and message for Gins JSON helpers
func ErrorMessage(err error) (code int, message map[string]interface{}) {
	v, ok := err.(*HTTPError)
	if ok {
		//code = v.Code()
		//if v.Code()/1000 == 6 {
		//	code = ErrInvalidParam.Code()
		//}
		return 200, map[string]interface{}{
			"type":        "error",
			"message":     v.Message(),
			"code":        v.Code(),
			"description": v.Desc(),
			"extra":       v.Extra(),
			"data":        v.Data(),
		}
	}

	return 200, map[string]interface{}{
		"message":     ErrInternalError.Error(),
		"code":        ErrInternalError.Code(),
		"description": err.Error(),
	}
}
