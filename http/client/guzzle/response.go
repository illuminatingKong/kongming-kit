package guzzle

import (
	"context"
	"net/http"
	"time"
)

type Response struct {
	Duration time.Duration

	// StatusCode is the HTTP Status Code returned by the HTTP Response. Taken from resp.StatusCode.
	StatusCode int

	// Header stores the response headers as http.Header interface.
	Header http.Header

	// Cookies stores the parsed response cookies.
	Cookies []*http.Cookie

	// Expose the native Go http.Response object for convenience.
	RawResponse *http.Response

	// Expose original request Context for convenience.
	Context *context.Context

	Err error
}

func (resp *Response) Json(useStruct interface{}) error {
	defer resp.Close()
	if resp.Err != nil {
		return resp.Err
	}

	if err := DecodeBody(resp.RawResponse, &useStruct); err != nil {
		return err
	}
	return nil
}

func (resp *Response) Xml(useStruct interface{}) error {
	defer resp.Close()
	if resp.Err != nil {
		return resp.Err
	}

	if err := XmlDecodeBody(resp.RawResponse, &useStruct); err != nil {
		return err
	}
	return nil
}
