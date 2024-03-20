package v5

import (
	"bytes"
	"errors"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/http/client/guzzle"
	"io"
)

type DefaultOption struct {
	DefaultPage     int
	DefaultPageSize int
}

type CMDBProvider struct {
	HttpClient *guzzle.Client
	AK         string
	SK         string
	Host       string
	Debug      bool
	Option     DefaultOption
	Logger     logx.Logger
}

func NewProvider(scheme, address, ak, sk string) *CMDBProvider {
	option := &guzzle.Config{
		Address:      address,
		Scheme:       scheme,
		ValidateHost: true,
	}
	httpClient, err := guzzle.NewClient(option)
	if err != nil {
		panic(err)
	}
	return &CMDBProvider{HttpClient: httpClient,
		AK: ak, SK: sk,
		Host:   address,
		Option: DefaultOption{DefaultPage: 1, DefaultPageSize: 20},
	}
}

func (c *CMDBProvider) OpenDebug() {
	c.Debug = true
}

func (c *CMDBProvider) Post(uri string, data interface{}, resp interface{}) (interface{}, error) {
	w := NewCMDBHttpRequesterClient(c.Host).WithAK(c.AK).WithSK(c.SK).WithUri(uri).WithBody(data)
	return c.postHandled(w, resp)
}

func (c *CMDBProvider) PostByte(uri string, data []byte, resp interface{}) (interface{}, error) {
	w := NewCMDBHttpRequesterClient(c.Host).WithAK(c.AK).WithSK(c.SK).WithUri(uri).WithBodyByte(data)
	return c.postHandled(w, resp)
}

func (c *CMDBProvider) postHandled(w *CMDBHttpRequester, resp interface{}) (interface{}, error) {
	w.WantPost()
	log := c.Logger
	if w.Err != nil {
		return nil, w.Err
	}
	r := c.HttpClient.DoNewRequest("POST", w.Uri)
	r.SetParam("Header", "Content-Type", "application/json")
	if len(w.HeaderHost) > 0 {
		r.SetParam("Header", "Host", w.HeaderHost)
	}

	r.SetParam("Header", "Accept", "*/*")
	r.SetParam("Query", "accesskey", w.AK)
	r.SetParam("Query", "signature", w.StringToSign)
	r.SetParam("Query", "expires", w.Expires)
	r.SetBody(w.Body)

	timeout, httpResp, err := c.HttpClient.NewDoRequest(r)
	if err != nil {
		return nil, err
	}
	out := guzzle.RequireOK(timeout, httpResp, err)
	if c.Debug == true {
		o, _ := io.ReadAll(out.RawResponse.Body)
		log.Infof("debug, post: %s, request: %+v , response: %+v", w.Uri, w.Body, string(o))
		out.RawResponse.Body = io.NopCloser(bytes.NewBuffer(o))
	}
	if out.StatusCode > 400 {
		rawBody := out.RawResponse.Body
		rb, _ := io.ReadAll(rawBody)
		nErr := errors.New(string(rb))
		_ = rawBody.Close()
		return nil, nErr
	}
	if err := out.Json(resp); err != nil {
		return nil, err
	}
	return resp, nil
}
